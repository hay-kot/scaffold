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
				EnvVars: []string{"SCAFFOLD_LOG_LEVEL"},
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

			//
			// Validate Runtime Config
			//
			err = rc.Validate()
			if err != nil {
				scaferrs := scaffoldrc.RcValidationErrors{}
				switch {
				case errors.As(err, &scaferrs):
					for _, err := range scaferrs {
						log.Error().Str("key", err.Key).Msg(err.Cause.Error())
					}
				default:
					return fmt.Errorf("unexpected error return from validator: %w", err)
				}
			}

			if rc.Settings.LogFile != "stdout" {
				f, err := os.Open(rc.Settings.LogFile)
				if err != nil {
					return fmt.Errorf("failed to open log file: %w", err)
				}

				log.Logger = log.Output(f)
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
			{
				Name:      "lint",
				Usage:     "lint a scaffoldrc file",
				UsageText: "scaffold lint [scaffold file]",
				Action:    ctrl.Lint,
			},
			{
				Name:      "init",
				Usage:     "initialize a new scaffold in the current directory for template scaffolds",
				UsageText: "scaffold init [flags]",
				Action:    ctrl.Init,
			},
			{
				Name: "dev",
				Subcommands: []*cli.Command{
					{
						Name:  "printer",
						Usage: "demos the printer",
						Action: func(ctx *cli.Context) error {
							console.Title(" --- Unknown Error ---")
							console.LineBreak()

							console.UnknownError("Basic Error", errors.New("this is a basic error's message"))

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
		case strings.HasPrefix(errstr, "flag provided but not defined"):
		default:
			log.Error().Err(err).Msg("error occurred")
			console.UnknownError("An unexpected error occurred", err)
		}

		os.Exit(1)
	}
}
