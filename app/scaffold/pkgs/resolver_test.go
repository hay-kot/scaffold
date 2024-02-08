package pkgs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestResolver_Resolve_Remote tests the Resolve method of the Resolver type with a remote
// url passed as an argument. We use a temporary directory to store the cloned repository.
// and override the Cloner with a ClonerFunc that creates the directory and returns the path.
//
// TODO: Update to test auth provider is called
func TestResolver_Resolve_Remote(t *testing.T) {
	tempdir := t.TempDir()

	tempcache := filepath.Join(tempdir, "cache")
	packagepath := filepath.Join(tempcache, "github.com", "hay-kot", "scaffold-go-cli")

	clonefn := ClonerFunc(func(path string, isBare bool, cfg *git.CloneOptions) (string, error) {
		t.Helper()

		err := os.MkdirAll(packagepath, 0755)
		require.NoError(t, err)

		return packagepath, nil
	})

	resolver := NewResolver(nil, tempcache, tempdir, WithCloner(clonefn))

	checkDirs := []string{tempdir}

	path, err := resolver.Resolve(
		"https://github.com/hay-kot/scaffold-go-cli",
		checkDirs,
		AuthProviderFunc(func(pkgurl string) (auth transport.AuthMethod, ok bool) {
			return nil, false
		}),
	)

	require.NoError(t, err)

	assert.Equal(t, packagepath, path)
}

func TestResolver_Resolve_Absolute(t *testing.T) {
}

func TestResolver_Resolve_Relative(t *testing.T) {
}

func TestResolver_Resolve_Cwd(t *testing.T) {
}
