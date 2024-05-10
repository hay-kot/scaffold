package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/urfave/cli/v2"
)

func (ctrl *Controller) Update(ctx *cli.Context) error {
	ctrl.ready()

	scaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return err
	}

	for _, s := range scaffolds {
		updated, err := pkgs.Update(filepath.Join(ctrl.Flags.Cache, s))
		if err != nil {
			return err
		}

		if updated {
			fmt.Printf("updated %s\n", s)
		}
	}

	fmt.Println("finished updating scaffolds")
	return nil
}
