package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/hay-kot/scaffold/app/commands"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

var (
	// Build information. Populated at build-time via -ldflags flag.
	version = "dev"
	commit  = "HEAD"
	date    = "now"
)

func build() string {
	short := commit
	if len(commit) > 7 {
		short = commit[:7]
	}

	return fmt.Sprintf("%s (%s) %s", version, short, date)
}

func HomeDir(s ...string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get home directory")
	}

	return filepath.Join(append([]string{home}, s...)...)
}

func main() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	ctrl := &commands.Controller{}

	app := &cli.App{
		Name:    "scaffold",
		Usage:   "scaffold projects and files from your terminal",
		Version: build(),
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "scaffoldrc",
				Usage:   "path to scaffoldrc file",
				Value:   HomeDir(".scaffold/scaffoldrc.yml"),
				EnvVars: []string{"SCAFFOLDRC"},
			},
			&cli.StringSliceFlag{
				Name:    "scaffold-dir",
				Usage:   "paths to directories containing scaffold templates",
				Value:   cli.NewStringSlice("./.scaffold"),
				EnvVars: []string{"SCAFFOLD_DIR"},
			},
			&cli.PathFlag{
				Name:    "cache",
				Usage:   "path to the local scaffold directory default",
				Value:   HomeDir(".scaffold/cache"),
				EnvVars: []string{"SCAFFOLD_CACHE"},
			},
			&cli.BoolFlag{
				Name:    "no-clobber",
				Usage:   "do not overwrite existing files (default: true)",
				EnvVars: []string{"SCAFFOLD_NO_CLOBBER"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "force",
				Usage:   "apply changes when git tree is dirty (default: false)",
				EnvVars: []string{"SCAFFOLD_FORCE"},
			},
			&cli.StringFlag{
				Name:    "output-dir",
				Usage:   "current working directory (where scaffold will be created)",
				Value:   ".",
				EnvVars: []string{"SCAFFOLD_OUT"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level (debug, info, warn, error, fatal, panic)",
				Value:   "warn",
				EnvVars: []string{"SCAFFOLD_LOG_LEVEL"},
			},
		},
		Before: func(ctx *cli.Context) error {
			ctrl.Flags = commands.Flags{
				NoClobber:      ctx.Bool("no-clobber"),
				Force:          ctx.Bool("force"),
				OutputDir:      ctx.String("output-dir"),
				Cache:          ctx.String("cache"),
				ScaffoldRCPath: ctx.String("scaffoldrc"),
				ScaffoldDirs:   ctx.StringSlice("scaffold-dir"),
			}

			switch ctx.String("log-level") {
			case "debug":
				log.Logger = log.Level(zerolog.DebugLevel)
			case "info":
				log.Logger = log.Level(zerolog.InfoLevel)
			case "warn":
				log.Logger = log.Level(zerolog.WarnLevel)
			case "error":
				log.Logger = log.Level(zerolog.ErrorLevel)
			case "fatal":
				log.Logger = log.Level(zerolog.FatalLevel)
			default:
				log.Logger = log.Level(zerolog.PanicLevel)
			}

			dir := filepath.Dir(ctrl.Flags.ScaffoldRCPath)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create scaffoldrc directory: %w", err)
			}

			if _, err := os.Stat(ctrl.Flags.ScaffoldRCPath); os.IsNotExist(err) {
				if err := os.WriteFile(ctrl.Flags.ScaffoldRCPath, []byte{}, 0o644); err != nil {
					return fmt.Errorf("failed to create scaffoldrc file: %w", err)
				}
			}

			if err := os.MkdirAll(ctx.String("cache"), 0o755); err != nil {
				return fmt.Errorf("failed to create cache directory: %w", err)
			}

			// Parse scaffoldrc file
			scaffoldrcFile, err := os.Open(ctrl.Flags.ScaffoldRCPath)
			if err != nil {
				return fmt.Errorf("failed to open scaffoldrc file: %w", err)
			}

			rc, err := scaffold.NewScaffoldRC(scaffoldrcFile)
			if err != nil {
				switch {
				case errors.As(err, &scaffold.RcValidationErrors{}):
					// I _know_ this is a valid cast, but the linter doesn't
					e := err.(scaffold.RcValidationErrors) //nolint:errorlint
					for _, err := range e {
						log.Error().Str("key", err.Key).Msg(err.Cause.Error())
					}
				default:
					return fmt.Errorf("failed to parse scaffoldrc file: %w", err)
				}
			}

			ctrl.Prepare(engine.New(), rc)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "new",
				Usage:     "create a new project from a scaffold",
				UsageText: "scaffold new [scaffold (url | path)] [flags]",
				Action:    ctrl.Project,
			},
			{
				Name:      "list",
				Usage:     "list available scaffolds",
				UsageText: "scaffold list [flags]",
				Action:    ctrl.List,
			},
			{
				Name:      "update",
				Usage:     "update the local cache of scaffolds",
				UsageText: "scaffold update [flags]",
				Action:    ctrl.Update,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run scaffold")
	}
}
