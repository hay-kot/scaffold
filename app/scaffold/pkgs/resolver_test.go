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

var noopAuthProvider = AuthProviderFunc(func(pkgurl string) (auth transport.AuthMethod, ok bool) {
	return nil, false
})

// TestResolver_Resolve_Remote tests the Resolve method of the Resolver type with a remote
// url passed as an argument. We use a temporary directory to store the cloned repository.
// and override the Cloner with a ClonerFunc that creates the directory and returns the path.
func TestResolver_Resolve_Remote(t *testing.T) {
	tempdir := t.TempDir()
	tempcache := filepath.Join(tempdir, "cache")
	packagepath := filepath.Join(tempcache, "github.com", "hay-kot", "scaffold-go-cli")

	clonefn := ClonerFunc(func(path string, version string, isBare bool, cfg *git.CloneOptions) (string, error) {
		t.Helper()

		err := os.MkdirAll(packagepath, 0o755)
		require.NoError(t, err)

		return packagepath, nil
	})

	resolver := NewResolver(nil, tempcache, tempdir, WithCloner(clonefn))

	path, err := resolver.Resolve(
		"https://github.com/hay-kot/scaffold-go-cli",
		nil,
		noopAuthProvider,
	)

	require.NoError(t, err)

	assert.Equal(t, packagepath, path)
}

// TestResolver_Resolve_Remote_Subdirectory tests the Resolve method of the Resolver type with
// a remote url passed as an argument. We use a temporary directory to store the cloned repository.
// and override the Cloner with a ClonerFunc that creates the directory and returns the path.
// In addition, we specify a subdirectory in the resolved url by appending a hash (#) and the
// subdirectory name. We expect the final path to contain the subdirectory too.

func TestResolver_Resolve_Remote_Subdirectory(t *testing.T) {
	tempdir := t.TempDir()
	tempcache := filepath.Join(tempdir, "cache")
	repopath := filepath.Join(tempcache, "github.com", "hay-kot", "scaffold-go-cli")
	subdirpath := filepath.Join(repopath, "subdir")

	clonefn := ClonerFunc(func(path string, version string, isBare bool, cfg *git.CloneOptions) (string, error) {
		t.Helper()

		err := os.MkdirAll(subdirpath, 0o755)
		require.NoError(t, err)

		// We always return the repo path, but the subdirectory must also be created
		return repopath, nil
	})

	resolver := NewResolver(nil, tempcache, tempdir, WithCloner(clonefn))

	path, err := resolver.Resolve(
		"https://github.com/hay-kot/scaffold-go-cli#subdir",
		nil,
		noopAuthProvider,
	)
	require.NoError(t, err)
	assert.Equal(t, subdirpath, path)
}

func TestResolver_Resolve_FilePaths(t *testing.T) {
	tempdir := t.TempDir()
	tempcwd := filepath.Join(tempdir, "cwd")
	tempcache := filepath.Join(tempdir, "cache")
	packagepath := filepath.Join(tempcwd, "scaffold", "test")

	// Setup .scaffold directory
	{
		err := os.MkdirAll(filepath.Join(tempcwd, ".scaffold", "new"), 0o755)
		require.NoError(t, err)

		// Create new/scaffold.yml file and new/templates dir

		f, err := os.Create(filepath.Join(tempcwd, ".scaffold", "new", "scaffold.yml"))
		require.NoError(t, err)
		_, _ = f.Write([]byte("---"))
		_ = f.Close()

		err = os.MkdirAll(filepath.Join(tempcwd, ".scaffold", "new", "templates"), 0o755)
		require.NoError(t, err)
	}

	err := os.MkdirAll(packagepath, 0o755)
	require.NoError(t, err)

	resolver := NewResolver(nil, tempcache, tempcwd)

	testcases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "relative path",
			input:    "scaffold/test",
			expected: packagepath,
		},
		{
			name:     "absolute path",
			input:    packagepath,
			expected: packagepath,
		},
		{
			name:     "local .scaffold directory",
			input:    "new",
			expected: filepath.Join(tempcwd, ".scaffold", "new"),
		},
	}

	checkDirs := []string{
		filepath.Join(tempcwd, ".scaffold"),
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			path, err := resolver.Resolve(tc.input, checkDirs, noopAuthProvider)

			require.NoError(t, err)

			assert.Equal(t, tc.expected, path)
		})
	}
}
