package pkgs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

type UpdateStatus int

const (
	UpdateStatusUpToDate UpdateStatus = iota
	UpdateStatusUpdated
	// UpdateStatusPinned means the repository was cloned at an @ref. A tag
	// or commit checkout is a detached HEAD that cannot pull, so Update
	// leaves pinned repositories as they are.
	UpdateStatusPinned
)

// Update pulls the latest commit for the cached repository that holds
// scaffoldPath. authprovider supplies credentials for the repository's origin
// URL.
func Update(cacheRoot string, scaffoldPath string, authprovider AuthProvider) (UpdateStatus, error) {
	repoPath, err := findRepoRoot(cacheRoot, scaffoldPath)
	if err != nil {
		return UpdateStatusUpToDate, err
	}

	if strings.Contains(filepath.Base(repoPath), "@") {
		return UpdateStatusPinned, nil
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return UpdateStatusUpToDate, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return UpdateStatusUpToDate, err
	}

	err = w.Pull(pullOptions(repo, authprovider))
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return UpdateStatusUpToDate, nil
		}
		return UpdateStatusUpToDate, err
	}

	return UpdateStatusUpdated, nil
}

// findRepoRoot walks up from scaffoldPath, because the scaffold may live in a
// nested directory of the repository.
func findRepoRoot(cacheRoot string, scaffoldPath string) (string, error) {
	parts := strings.Split(filepath.ToSlash(scaffoldPath), "/")

	// A cached repository has at least two parts, i.e. git.host.com/repo.
	for i := len(parts); i > 1; i-- {
		repoPath := filepath.Join(cacheRoot, filepath.Join(parts[:i]...))
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
			return repoPath, nil
		}
	}

	return "", fmt.Errorf("no git repository found for %s", scaffoldPath)
}

func pullOptions(repo *git.Repository, authprovider AuthProvider) *git.PullOptions {
	opts := &git.PullOptions{}

	remote, err := repo.Remote(git.DefaultRemoteName)
	if err != nil || len(remote.Config().URLs) == 0 {
		return opts
	}

	if auth, ok := authprovider.Authenticator(remote.Config().URLs[0]); ok {
		opts.Auth = auth
	}

	return opts
}
