package commands

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSyncRCScaffolds_WarnsAndContinues(t *testing.T) {
	cache := t.TempDir()
	cachedDir := filepath.Join(cache, "github.com", "org", "cached")
	require.NoError(t, os.MkdirAll(cachedDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(cachedDir, "scaffold.yaml"), nil, 0o644))

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

func TestSyncRCScaffolds_RedactsCredentials(t *testing.T) {
	rc := scaffoldrc.Default()
	rc.Scaffolds = []string{"https://bot:ghp_SECRET123@127.0.0.1:1/org/private"}

	ctrl := &Controller{Flags: Flags{Cache: t.TempDir()}}
	ctrl.Prepare(engine.New(), rc)

	stderr := &bytes.Buffer{}
	ctrl.errPrinter = newPrinter(stderr)

	ctrl.syncRCScaffolds()

	assert.Contains(t, stderr.String(), "https://redacted@127.0.0.1:1/org/private")
	assert.NotContains(t, stderr.String(), "ghp_SECRET123")
}

func TestSyncRCScaffolds_KeepsCloneProgressOffStdout(t *testing.T) {
	rc := scaffoldrc.Default()
	rc.Scaffolds = []string{"https://github.com/org/repo"}

	var clones int
	ctrl := &Controller{Flags: Flags{Cache: t.TempDir()}}
	ctrl.Prepare(engine.New(), rc)
	ctrl.syncResolverOpts = []pkgs.ResolverOption{pkgs.WithCloner(pkgs.ClonerFunc(
		func(ctx context.Context, path string, version string, isBare bool, cfg *git.CloneOptions) (string, error) {
			clones++
			assert.NotEqual(t, os.Stdout, cfg.Progress, "clone progress must not reach stdout")
			return path, nil
		},
	))}

	ctrl.syncRCScaffolds()
	assert.Equal(t, 1, clones)
}
