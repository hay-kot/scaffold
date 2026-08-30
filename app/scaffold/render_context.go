package scaffold

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/rs/zerolog/log"
)

// memoryOutputDir is the sentinel the CLI uses to render into an in-memory
// filesystem. It is not a real path, so .Ctx reports empty output paths for it.
const memoryOutputDir = ":memory:"

// buildCtxVars collects the filesystem and repository facts a template may need
// to derive values that differ per checkout, such as a port block or a docker
// stack name. Every field is best effort: a missing git repository or an
// unreadable working directory yields an empty string rather than an error,
// because none of this is required to render a scaffold.
func buildCtxVars(outputDir string) engine.Vars {
	workingDir, err := os.Getwd()
	if err != nil {
		log.Debug().Err(err).Msg("failed to resolve working directory for .Ctx")
	}

	var outputDirAbs string
	if outputDir != "" && outputDir != memoryOutputDir {
		outputDirAbs, err = filepath.Abs(outputDir)
		if err != nil {
			log.Debug().Err(err).Str("dir", outputDir).Msg("failed to resolve absolute output directory for .Ctx")
			outputDirAbs = ""
		}
	}

	// The basename of "." or "./" is not useful as a seed, so resolve it through
	// the absolute path. This is the common case: scaffold runs in the directory
	// it renders into.
	outputDirBase := ""
	if outputDirAbs != "" {
		outputDirBase = filepath.Base(outputDirAbs)
	}

	workingDirBase := ""
	if workingDir != "" {
		workingDirBase = filepath.Base(workingDir)
	}

	gitDir := outputDirAbs
	if gitDir == "" {
		gitDir = workingDir
	}

	return engine.Vars{
		"OutputDir":      outputDir,
		"OutputDirAbs":   outputDirAbs,
		"OutputDirBase":  outputDirBase,
		"WorkingDir":     workingDir,
		"WorkingDirBase": workingDirBase,
		"Git":            buildGitVars(gitDir),
	}
}

// buildGitVars reads the branch and repository name from the git repository that
// contains dir. A detached HEAD reports an empty branch, since there is no name
// to report.
func buildGitVars(dir string) engine.Vars {
	out := engine.Vars{
		"Branch": "",
		"Repo":   "",
	}

	if dir == "" {
		return out
	}

	repo, err := git.PlainOpenWithOptions(dir, &git.PlainOpenOptions{DetectDotGit: true})
	if err != nil {
		log.Debug().Err(err).Str("dir", dir).Msg("no git repository for .Ctx.Git")
		return out
	}

	head, err := repo.Head()
	if err != nil {
		log.Debug().Err(err).Msg("failed to read git HEAD for .Ctx.Git")
	} else if head.Name().IsBranch() {
		out["Branch"] = head.Name().Short()
	}

	// go-git has no notion of a repository name, so derive it from the worktree
	// root. That matches what a user sees on disk, which is what a seed wants.
	if wt, err := repo.Worktree(); err == nil {
		out["Repo"] = filepath.Base(wt.Filesystem.Root())
	} else {
		log.Debug().Err(err).Msg("failed to read git worktree for .Ctx.Git")
	}

	return out
}
