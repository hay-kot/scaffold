package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hay-kot/scaffold/app/commands"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"github.com/hay-kot/scaffold/internal/printer"
	"github.com/hay-kot/scaffold/internal/styles"
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

var ErrLinterErrors = errors.New("scaffold errors found")

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
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).Level(zerolog.WarnLevel)

	ctrl := &commands.Controller{}
	console := printer.New(os.Stdout)

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
				Usage:   "do not overwrite existing files",
				EnvVars: []string{"SCAFFOLD_NO_CLOBBER"},
				Value:   true,
			},
			&cli.BoolFlag{
				Name:    "force",
				Usage:   "apply changes when git tree is dirty",
				Value:   true,
				EnvVars: []string{"SCAFFOLD_FORCE"},
			},
			&cli.StringFlag{
				Name:    "output-dir",
				Usage:   "scaffold output directory (use ':memory:' for in-memory filesystem)",
				Value:   ".",
				EnvVars: []string{"SCAFFOLD_OUT"},
			},
			&cli.StringFlag{
				Name:    "log-level",
				Usage:   "log level (debug, info, warn, error, fatal, panic)",
				Value:   "warn",
				EnvVars: []string{"SCAFFOLD_LOG_LEVEL", "SCAFFOLD_SETTINGS_LOG_LEVEL"},
			},
			&cli.StringFlag{
				Name:    "log-file",
				Usage:   "log file to write to (use 'stdout' for stdout)",
				EnvVars: []string{"SCAFFOLD_SETTINGS_LOG_FILE"},
			},
			&cli.StringFlag{
				Name:    "theme",
				Usage:   "theme to use for the scaffold output",
				Value:   "scaffold",
				EnvVars: []string{"SCAFFOLD_THEME", "SCAFFOLD_SETTINGS_THEME"},
			},
			&cli.StringFlag{
				Name:    "run-hooks",
				Usage:   "run hooks (never, always, prompt) when provided overrides scaffold rc",
				EnvVars: []string{"SCAFFOLD_SETTINGS_RUN_HOOKS"},
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
				if !errors.Is(err, os.ErrNotExist) {
					return fmt.Errorf("failed to open scaffoldrc file: %w", err)
				}
				log.Debug().Msg("scaffoldrc file does not exist, skipping")
			}

			rc := scaffoldrc.Default()

			if scaffoldrcFile != nil {
				rc, err = scaffoldrc.New(scaffoldrcFile)
				if err != nil {
					return err
				}
			}

			//
			// Override Settings with Flags
			//
			if ctx.IsSet("theme") {
				rc.Settings.Theme = styles.HuhTheme(ctx.String("theme"))
			}

			if ctx.IsSet("run-hooks") {
				rc.Settings.RunHooks = scaffoldrc.ParseRunHooksOption(ctx.String("run-hooks"))
			}

			if ctx.IsSet("log-level") {
				level, err := zerolog.ParseLevel(ctx.String("log-level"))
				if err != nil {
					return fmt.Errorf("failed to parse log level: %w", err)
				}

				log.Logger = log.Level(level)
			}

			if ctx.IsSet("log-file") {
				rc.Settings.LogFile = ctx.String("log-file")

				if !strings.HasPrefix(rc.Settings.LogFile, "/") {
					// If the file path is not absolute, we want to make it absolute
					// so that it is relative to the cwd and not the scaffoldrc file.
					absLogFilePath, err := filepath.Abs(rc.Settings.LogFile)
					if err != nil {
						return err
					}

					rc.Settings.LogFile = absLogFilePath
				}
			}

			//
			// Validate Runtime Config
			//
			err = rc.Validate()
			if err != nil {
				scaferrs := scaffoldrc.RcValidationErrors{}
				switch {
				case errors.As(err, &scaferrs):
					errlist := make([]printer.KeyValueError, 0, len(scaferrs))
					for _, err := range scaferrs {
						errlist = append(errlist, printer.KeyValueError{Key: err.Key, Message: err.Cause.Error()})
					}

					console.KeyValueValidationError("ScaffoldRC Errors", errlist)
				default:
					return fmt.Errorf("unexpected error return from validator: %w", err)
				}
			}

			if rc.Settings.LogFile != "stdout" {
				logpath := rc.Settings.LogFile
				if !strings.HasPrefix(ctrl.Flags.ScaffoldRCPath, "/") {
					// Assume that the path is relative to the scaffold rc file
					logpath = filepath.Join(dir, rc.Settings.LogFile)
				}

				f, err := os.OpenFile(logpath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
				if err != nil {
					return fmt.Errorf("failed to open log file: %w", err)
				}

				log.Logger = log.Output(zerolog.ConsoleWriter{
					Out:     f,
					NoColor: true,
				})
			}

			styles.SetGlobalStyles(rc.Settings.Theme)
			console = console.WithBase(styles.Base).WithLight(styles.Light)

			ctrl.Prepare(engine.New(), rc)
			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "new",
				Usage:     "create a new project from a scaffold",
				UsageText: "scaffold new [scaffold (url | path)] [flags]",
				Flags: []cli.Flag{
					&cli.BoolFlag{
						Name:  "no-prompt",
						Usage: "disable interactive mode",
						Value: false,
					},
					&cli.StringFlag{
						Name:  "preset",
						Usage: "preset to use for the scaffold",
						Value: "",
					},
					&cli.StringFlag{
						Name:  "snapshot",
						Usage: "path or `stdout` to save the output ast",
						Value: "",
					},
				},
				Action: func(ctx *cli.Context) error {
					return ctrl.New(ctx.Args().Slice(), commands.FlagsNew{
						NoPrompt: ctx.Bool("no-prompt"),
						Preset:   ctx.String("preset"),
						Snapshot: ctx.String("snapshot"),
					})
				},
			},
			{
				Name:   "list",
				Usage:  "list available scaffolds",
				Action: ctrl.List,
			},
			{
				Name:   "update",
				Usage:  "update the local cache of scaffolds",
				Action: ctrl.Update,
			},
			{
				Name:      "lint",
				Usage:     "lint a scaffoldrc file",
				UsageText: "scaffold lint [scaffold file]",
				Action: func(ctx *cli.Context) error {
					pfpath := ctx.Args().First()
					if pfpath == "" {
						return errors.New("no file provided")
					}

					err := ctrl.Lint(pfpath)
					if err != nil {
						errlist, ok := err.(commands.ErrList) // nolint: errorlint
						if !ok {
							return err
						}

						items := make([]printer.StatusListItem, 0, len(errlist))
						for _, e := range errlist {
							items = append(items, printer.StatusListItem{Ok: false, Status: e.Error()})
						}

						console.StatusList("Scaffold Errors", items)

						return ErrLinterErrors
					}

					return nil
				},
			},
			{
				Name:   "init",
				Usage:  "initialize a new scaffold in the current directory for template scaffolds",
				Action: ctrl.Init,
			},
			{
				Name:  "dev",
				Usage: "development commands for testing",
				Subcommands: []*cli.Command{
					{
						Name:  "printer",
						Usage: "demos the printer",
						Action: func(ctx *cli.Context) error {
							console.Title(" --- Unknown Error ---")
							console.LineBreak()

							console.FatalError(errors.New("this is a basic error's message"))

							console.LineBreak()
							console.Title(" --- List ---")
							console.LineBreak()

							console.List("List Items", []string{"item 1", "item 2", "item 3"})

							console.LineBreak()
							console.Title(" --- StatusList ---")
							console.LineBreak()

							console.StatusList("Status Items", []printer.StatusListItem{
								{Ok: true, Status: "Status 1"},
								{Ok: false, Status: "Status 2"},
								{Ok: true, Status: "Status 3"},
							})

							console.LineBreak()
							console.Title(" --- Key Value Error ---")
							console.LineBreak()

							console.KeyValueValidationError("Key Value Errors", []printer.KeyValueError{
								{Key: "alias.gh", Message: "invalid choice for key_1"},
								{Key: "settings.theme", Message: "invalid theme 'x-theme'"},
							})

							return nil
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		errstr := err.Error()

		switch {
		// ignore these errors, urfave/cli does not provide any way to hanldle them
		// without direct string comparison :(
		case strings.HasPrefix(errstr, "flag provided but not defined"), errors.Is(err, ErrLinterErrors):
			// ignore
		default:
			console.FatalError(err)
		}

		os.Exit(1)
	}
}
