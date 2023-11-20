package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/urfave/cli/v2"
)

func (ctrl *Controller) List(ctx *cli.Context) error {
	systemScaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return err
	}

	localScaffolds, err := pkgs.ListLocal(os.DirFS(ctrl.Flags.OutputDir))
	if err != nil {
		return err
	}

	bldr := strings.Builder{}

	if len(localScaffolds) > 0 {
		bldr.WriteString("## Local Scaffolds\n")

		for _, s := range localScaffolds {
			bldr.WriteString(fmt.Sprintf("  - %s\n", s))
		}

		bldr.WriteString("\n")
	}

	if len(systemScaffolds) > 0 {
		bldr.WriteString("## System Scaffolds\n")

		for _, s := range systemScaffolds {
			bldr.WriteString(fmt.Sprintf("  - %s\n", s))
		}
	}

	out, err := glamour.RenderWithEnvironmentConfig(bldr.String())
	if err != nil {
		return err
	}

	fmt.Println(out)

	return nil
}
