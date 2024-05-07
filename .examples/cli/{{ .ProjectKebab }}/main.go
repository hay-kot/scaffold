package main

import (
	"fmt"
	"os"

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
	app := &cli.App{
		Name:    "{{ .Project }}",
		Usage:   "{{ .Scaffold.description }}",
		Version: build(),
		Commands: []*cli.Command{
			{
				Name:  "hello",
				Usage: "Says hello world",
				Action: func(ctx *cli.Context) error {
					fmt.Println("Hello, your favorite colors are {{ .Scaffold.colors | join `, ` }}")
					return nil
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Err(err).Msg("failed to run scaffold")
	}
}
