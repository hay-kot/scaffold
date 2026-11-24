package pkgs

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/object"
	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func commitFile(t *testing.T, repo *git.Repository, dir, name string) {
	t.Helper()

	require.NoError(t, os.WriteFile(filepath.Join(dir, name), []byte(name), 0o644))

	wt, err := repo.Worktree()
	require.NoError(t, err)

	_, err = wt.Add(name)
	require.NoError(t, err)

	_, err = wt.Commit("add "+name, &git.CommitOptions{
		Author: &object.Signature{Name: "test", Email: "test@example.com", When: time.Now()},
	})
	require.NoError(t, err)
}

func TestUpdate(t *testing.T) {
	tempdir := t.TempDir()
	upstreamDir := filepath.Join(tempdir, "upstream")
	cache := filepath.Join(tempdir, "cache")
	cloneDir := filepath.Join(cache, "example.com", "org", "repo")

	upstream, err := git.PlainInit(upstreamDir, false)
	require.NoError(t, err)
	commitFile(t, upstream, upstreamDir, "one.txt")

	_, err = git.PlainClone(cloneDir, false, &git.CloneOptions{URL: upstreamDir})
	require.NoError(t, err)

	var authURL string
	recordAuth := AuthProviderFunc(func(pkgurl string) (transport.AuthMethod, bool) {
		authURL = pkgurl
		return nil, false
	})

	status, err := Update(context.Background(), cache, "example.com/org/repo", recordAuth)
	require.NoError(t, err)
	assert.Equal(t, UpdateStatusUpToDate, status)
	assert.Equal(t, upstreamDir, authURL, "auth lookup must use the origin URL")

	commitFile(t, upstream, upstreamDir, "two.txt")

	status, err = Update(context.Background(), cache, "example.com/org/repo", recordAuth)
	require.NoError(t, err)
	assert.Equal(t, UpdateStatusUpdated, status)
	assert.FileExists(t, filepath.Join(cloneDir, "two.txt"))
}

func TestUpdate_Pinned(t *testing.T) {
	cache := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(cache, "example.com", "org", "repo@v1.0.0", ".git"), 0o755))

	status, err := Update(context.Background(), cache, "example.com/org/repo@v1.0.0", noopAuthProvider)
	require.NoError(t, err)
	assert.Equal(t, UpdateStatusPinned, status)
}

func TestUpdate_NoRepository(t *testing.T) {
	_, err := Update(context.Background(), t.TempDir(), "example.com/org/repo", noopAuthProvider)
	require.Error(t, err)
}

func TestPullOptions_UsesOriginAuth(t *testing.T) {
	repo, err := git.PlainInit(t.TempDir(), false)
	require.NoError(t, err)

	_, err = repo.CreateRemote(&config.RemoteConfig{Name: git.DefaultRemoteName, URLs: []string{"https://example.com/org/repo"}})
	require.NoError(t, err)

	want := &githttp.TokenAuth{Token: "sentinel"}
	provider := AuthProviderFunc(func(pkgurl string) (transport.AuthMethod, bool) {
		if pkgurl == "https://example.com/org/repo" {
			return want, true
		}
		return nil, false
	})

	assert.Same(t, want, pullOptions(repo, provider).Auth)
	assert.Nil(t, pullOptions(repo, noopAuthProvider).Auth)
}
