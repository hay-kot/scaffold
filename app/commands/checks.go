package commands

import (
	"errors"

	"github.com/go-git/go-git/v5"
	"github.com/rs/zerolog/log"
)

func checkWorkingTree(dir string) bool {
	log.Debug().Str("dir", dir).Msg("opening git repository")
	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Debug().Err(err).Msg("failed to open git repository")
		return errors.Is(err, git.ErrRepositoryNotExists)
	}

	log.Debug().Msg("opening git worktree")
	wt, err := repo.Worktree()
	if err != nil {
		log.Debug().Err(err).Msg("failed to open git worktree")
		return false
	}

	log.Debug().Msg("checking git status")
	status, err := wt.Status()
	if err != nil {
		log.Debug().Err(err).Msg("failed to get git status")
		return false
	}

	if status.IsClean() {
		return true
	}

	return false
}
