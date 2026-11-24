package commands

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncRCScaffolds_WarnsAndContinues(t *testing.T) {
	cache := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(cache, "github.com", "org", "cached"), 0o755))

	rc := scaffoldrc.Default()
	rc.Scaffolds = []string{
		"https://github.com/org/cached",
		"https://127.0.0.1:1/org/unreachable",
	}

	ctrl := &Controller{Flags: Flags{Cache: cache}}
	ctrl.Prepare(engine.New(), rc)

	stdout, stderr := &bytes.Buffer{}, &bytes.Buffer{}
	ctrl.printer = newPrinter(stdout)
	ctrl.errPrinter = newPrinter(stderr)

	ctrl.syncRCScaffolds()

	assert.Empty(t, stdout.String(), "sync must not write to stdout")
	assert.Contains(t, stderr.String(), "https://127.0.0.1:1/org/unreachable")
	assert.NotContains(t, stderr.String(), "org/cached")
	assert.NoDirExists(t, filepath.Join(cache, "127.0.0.1:1", "org", "unreachable"))
}
