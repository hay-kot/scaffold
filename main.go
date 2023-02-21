package main

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/hay-kot/scaffold/internal/core/rwfs"
	"github.com/hay-kot/scaffold/internal/engine"
	"github.com/hay-kot/scaffold/scaffold"
	"github.com/hay-kot/scaffold/scaffold/pkgs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"

	"github.com/go-git/go-git/v5"
)

// TODO (hay-kot): add support for remote repositories with specific tags
// Users should be able to provide url with @tag to scaffold a specific tag
// from a remote repository.

// TODO (hay-kot): check for changes on remote repository and pull if needed
// When using a remote repository, scaffold should check for changes on the
// remote repository and pull if needed. Allow for disabling this feature via
// the CLI `--pull false` flag.

// TODO (hay-kot): merge defaults from scaffoldrc and command line args.
// The scaffoldrc file should be merged with the command line arguments to
// provide a single configuration for the scaffold command. allowing users to
// override the default configuration.

// TODO (hay-kot): remove --vars and use args instead
// User should be able to append key=value pairs to the scaffold command to
// provide variables for the template. This will allow users to provide
// variables without having to specify a flag.

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

	ctrl := &controller{}

	app := &cli.App{
		Name:    "scaffold",
		Usage:   "scaffold projects and files from your terminal",
		Version: build(),
		Flags: []cli.Flag{
			&cli.PathFlag{
				Name:    "scaffoldrc",
				Usage:   "path to scaffoldrc file",
				Value:   HomeDir(".scaffold/scaffoldrc.yml"),
				EnvVars: []string{"SCAFFOLD_RC"},
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
				Name:    "out",
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
			&cli.StringSliceFlag{
				Name:  "var",
				Usage: "key/value pairs to use as variables in the scaffold (e.g. --var foo=bar)",
			},
		},
		Before: func(ctx *cli.Context) error {
			ctrl.logLevel = ctx.String("log-level")
			switch ctrl.logLevel {
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

			// Creates scaffoldrc file if it doesn't exist
			scaffoldrc := ctx.String("scaffoldrc")

			dir := filepath.Dir(scaffoldrc)
			if err := os.MkdirAll(dir, 0o755); err != nil {
				return fmt.Errorf("failed to create scaffoldrc directory: %w", err)
			}

			if _, err := os.Stat(scaffoldrc); os.IsNotExist(err) {
				if err := os.WriteFile(scaffoldrc, []byte{}, 0o644); err != nil {
					return fmt.Errorf("failed to create scaffoldrc file: %w", err)
				}
			}

			if err := os.MkdirAll(ctx.String("cache"), 0o755); err != nil {
				return fmt.Errorf("failed to create cache directory: %w", err)
			}

			// Read Global Flags
			ctrl.cwd = ctx.String("out")

			ctrl.noClobber = ctx.Bool("no-clobber")
			ctrl.scaffoldrc = scaffoldrc
			ctrl.cache = ctx.String("cache")

			varString := ctx.StringSlice("var")
			ctrl.vars = make(map[string]string, len(varString))
			for _, v := range varString {
				kv := strings.Split(v, "=")
				ctrl.vars[kv[0]] = kv[1]
			}

			// Parse scaffoldrc file
			scaffoldrcFile, err := os.Open(ctrl.scaffoldrc)
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

			ctrl.rc = rc
			ctrl.engine = engine.New()

			return nil
		},
		Commands: []*cli.Command{
			{
				Name:      "new",
				Usage:     "create a new project from a scaffold",
				UsageText: "scaffold new [scaffold (url | path)] [flags]",
				Action:    ctrl.Project,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run scaffold")
	}
}

type controller struct {
	engine *engine.Engine

	// Global Flags

	rc         *scaffold.ScaffoldRC
	cwd        string
	vars       map[string]string
	logLevel   string
	noClobber  bool
	scaffoldrc string
	cache      string
	force      bool
}

func (c *controller) Project(ctx *cli.Context) error {
	argPath := ctx.Args().First()
	if argPath == "" {
		return fmt.Errorf("path is required")
	}

	if !c.force {
		ok := checkWorkingTree(c.cwd)
		if !ok {
			log.Warn().Msg("working tree is dirty, use --force to apply changes")
			return nil
		}
	}

	resolver := pkgs.NewResolver(
		c.rc.Shorts,
		c.cache,
		".",
	)

	path, err := resolver.Resolve(argPath)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	pfs := os.DirFS(path)

	p, err := scaffold.LoadProject(pfs, scaffold.Options{
		NoClobber: c.noClobber,
	})
	if err != nil {
		return err
	}

	vars, err := p.AskQuestions(c.vars)
	if err != nil {
		return err
	}

	args := &scaffold.RWFSArgs{
		Project: p,
		ReadFS:  pfs,
		WriteFS: rwfs.NewOsWFS(c.cwd),
	}

	err = scaffold.RenderRWFS(c.engine, args, vars)

	if err != nil {
		return err
	}

	return nil
}

func checkWorkingTree(dir string) bool {
	repo, err := git.PlainOpen(dir)
	if err != nil {
		log.Debug().Err(err).Msg("failed to open git repository")
		return errors.Is(err, git.ErrRepositoryNotExists)
	}

	wt, err := repo.Worktree()
	if err != nil {
		log.Debug().Err(err).Msg("failed to open git worktree")
		return false
	}

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
