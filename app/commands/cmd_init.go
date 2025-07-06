package commands

import (
	"context"
	"embed"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/urfave/cli/v3"
)

//go:embed init/*
var initFiles embed.FS

func (ctrl *Controller) Init(_ context.Context, c *cli.Command) error {
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

	return nil
}
