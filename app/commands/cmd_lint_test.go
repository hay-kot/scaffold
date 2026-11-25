package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func writeLintScaffold(t *testing.T, dirs ...string) string {
	t.Helper()

	root := t.TempDir()
	pfpath := filepath.Join(root, "scaffold.yaml")
	require.NoError(t, os.WriteFile(pfpath, []byte("questions: []\n"), 0o644))

	for _, dir := range dirs {
		require.NoError(t, os.MkdirAll(filepath.Join(root, dir), 0o755))
	}

	return pfpath
}

func TestLint_RootDir(t *testing.T) {
	ctrl := &Controller{}

	require.NoError(t, ctrl.Lint(writeLintScaffold(t, "templates")))
	require.ErrorContains(t, ctrl.Lint(writeLintScaffold(t)), "does not exist")
	require.ErrorContains(t, ctrl.Lint(writeLintScaffold(t, "templates", "{{ .Project }}")), "exactly one root directory")
}
