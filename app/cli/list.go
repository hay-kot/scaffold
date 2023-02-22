package handlers

import (
	"fmt"
	"os"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/urfave/cli/v2"
)

func (ctrl *Controller) List(ctx *cli.Context) error {
	scaffolds, err := pkgs.List(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return err
	}

	for _, s := range scaffolds {
		fmt.Println(s)
	}
	return nil
}
