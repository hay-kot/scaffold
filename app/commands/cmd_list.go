package commands

import (
	"os"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
)

type FlagsList struct {
	OutputDir string
}

func (ctrl *Controller) List(flags FlagsList) error {
	systemScaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return err
	}

	localScaffolds, err := pkgs.ListLocal(os.DirFS(flags.OutputDir))
	if err != nil {
		return err
	}

	ctrl.printer.LineBreak()

	if len(localScaffolds) > 0 {
		ctrl.printer.List("Local Scaffolds", localScaffolds)
	}

	if len(systemScaffolds) > 0 {
		ctrl.printer.List("System Scaffolds", systemScaffolds)
	}

	ctrl.printer.LineBreak()
	return nil
}
