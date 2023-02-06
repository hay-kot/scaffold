package main

import (
	"fmt"
	"strings"

	"os"

	"github.com/hay-kot/scaffold/internal/core/rwfs"
	"github.com/hay-kot/scaffold/internal/engine"
	"github.com/hay-kot/scaffold/scaffold"
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

func main() {
	ctrl := &controller{}

	app := &cli.App{
		Name:    "scaffold",
		Usage:   "scaffold projects and files from your terminal",
		Version: build(),
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "cwd",
				Usage: "current working directory (where scaffold will be created)",
				Value: ".",
			},
			&cli.StringFlag{
				Name:  "log-level",
				Usage: "log level (debug, info, warn, error, fatal, panic)",
				Value: "panic",
			},
			&cli.StringSliceFlag{
				Name:  "var",
				Usage: "key/value pairs to use as variables in the scaffold (e.g. --var foo=bar)",
			},
		},
		Before: func(ctx *cli.Context) error {
			ctrl.engine = engine.New()

			ctrl.cwd = ctx.String("cwd")
			ctrl.logLevel = ctx.String("log-level")

			varString := ctx.StringSlice("var")

			ctrl.vars = make(map[string]string, len(varString))
			for _, v := range varString {
				kv := strings.Split(v, "=")
				ctrl.vars[kv[0]] = kv[1]
			}

			log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
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
				Name:      "scaffold",
				Usage:     "create a new scaffold",
				UsageText: "scaffold scaffold [name] [flags]",
				Action:    ctrl.Scaffold,
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run scaffold")
	}
}

type controller struct {
	cwd      string
	vars     map[string]string
	logLevel string
	engine   *engine.Engine
}

func (c *controller) Scaffold(ctx *cli.Context) error {
	log.Panic().Msg("not implemented")
	return nil
}

func (c *controller) Project(ctx *cli.Context) error {
	path := ctx.Args().First()
	if path == "" {
		return fmt.Errorf("path is required")
	}

	pfs := os.DirFS(path)

	p, err := scaffold.LoadProject(pfs)
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
