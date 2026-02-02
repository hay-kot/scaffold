package commands

import (
	"context"
	"embed"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/urfave/cli/v3"
)

//go:embed init/*
var initFiles embed.FS

// FlagsInit contains flags for the init command
type FlagsInit struct {
	Stealth bool
}

func (ctrl *Controller) Init(_ context.Context, c *cli.Command) error {
	flags := FlagsInit{
		Stealth: c.Bool("stealth"),
	}

	// get current directory
	dir, err := os.Getwd()
	if err != nil {
		return err
	}

	scaffoldDir := filepath.Join(dir, ".scaffold")

	// Check if exists
	if _, err := os.Stat(scaffoldDir); !os.IsNotExist(err) {
		return err
	}

	// Create directory
	if err := os.Mkdir(scaffoldDir, 0o755); err != nil {
		return err
	}

	// Write files from initFiles embed.FS to .scaffold
	files, err := doublestar.Glob(initFiles, "init/**/*")
	if err != nil {
		return err
	}

	for _, file := range files {
		f, err := initFiles.Open(file)
		if err != nil {
			return err
		}

		// Check if dir
		finfo, err := f.Stat()
		if err != nil {
			return err
		}

		fileName := strings.TrimPrefix(file, "init/")

		if finfo.IsDir() {
			err := os.Mkdir(filepath.Join(scaffoldDir, fileName), 0o755)
			if err != nil {
				return err
			}

			continue
		}

		out, err := os.Create(filepath.Join(scaffoldDir, fileName))
		if err != nil {
			return err
		}

		if _, err := io.Copy(out, f); err != nil {
			return err
		}

		if err := out.Close(); err != nil {
			return err
		}

		if err := f.Close(); err != nil {
			return err
		}
	}

	// Handle stealth mode - add .scaffold to .git/info/exclude
	if flags.Stealth {
		if err := addToGitExclude(dir, ".scaffold"); err != nil {
			ctrl.printer.Warning(fmt.Sprintf("Warning: could not add .scaffold to git exclude: %v", err))
		}
	}

	return nil
}

// addToGitExclude adds an entry to .git/info/exclude if not already present
func addToGitExclude(repoDir, entry string) error {
	gitDir := filepath.Join(repoDir, ".git")

	// Check if .git directory exists
	info, err := os.Stat(gitDir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("not a git repository")
		}
		return err
	}

	// Handle git worktrees - .git might be a file pointing to the actual git dir
	if !info.IsDir() {
		content, err := os.ReadFile(gitDir)
		if err != nil {
			return err
		}
		// Parse "gitdir: /path/to/actual/git/dir"
		gitdirLine := strings.TrimSpace(string(content))
		if strings.HasPrefix(gitdirLine, "gitdir: ") {
			gitDir = strings.TrimPrefix(gitdirLine, "gitdir: ")
		}
	}

	excludePath := filepath.Join(gitDir, "info", "exclude")

	// Ensure .git/info directory exists
	infoDir := filepath.Dir(excludePath)
	if err := os.MkdirAll(infoDir, 0o755); err != nil {
		return err
	}

	// Read existing exclude file if it exists
	var existingContent []byte
	if _, err := os.Stat(excludePath); err == nil {
		existingContent, err = os.ReadFile(excludePath)
		if err != nil {
			return err
		}
	}

	// Check if entry already exists
	lines := strings.Split(string(existingContent), "\n")
	for _, line := range lines {
		if strings.TrimSpace(line) == entry {
			return nil // Already excluded
		}
	}

	// Append entry to exclude file
	f, err := os.OpenFile(excludePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck

	// Add newline before entry if file doesn't end with one
	prefix := ""
	if len(existingContent) > 0 && !strings.HasSuffix(string(existingContent), "\n") {
		prefix = "\n"
	}

	_, err = f.WriteString(prefix + entry + "\n")
	return err
}
