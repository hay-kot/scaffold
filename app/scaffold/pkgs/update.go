package pkgs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// Update updates a git repository to the latest commit
func Update(cacheRoot string, scaffoldPath string) (updated bool, err error) {
	// The package may be in a nested directory of the repository, so we need to find the root
	// of the repository

	scaffoldPathParts := strings.Split(filepath.ToSlash(scaffoldPath), "/")

	repoPath := ""
	gitDirFound := false

	// Assuming a minimum of 2 parts in the path i.e. git.host.com/repo
	for i := len(scaffoldPathParts); i > 1; i-- {
		repoPath = filepath.Join(cacheRoot, filepath.Join(scaffoldPathParts[:i]...))
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
			gitDirFound = true
			break
		}
	}

	if !gitDirFound {
		return false, fmt.Errorf("no git repository found for %s", scaffoldPath)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return false, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	err = w.Pull(&git.PullOptions{})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
