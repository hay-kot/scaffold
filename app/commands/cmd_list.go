package commands

import (
	"os"
	"slices"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/internal/printer"
)

type FlagsList struct {
	OutputDir string
}

func (ctrl *Controller) ListPlain(flags FlagsList) ([]string, error) {
	systemScaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return nil, err
	}

	localScaffolds, err := pkgs.ListLocal(os.DirFS(flags.OutputDir))
	if err != nil {
		return nil, err
	}

	out := make([]string, 0, len(systemScaffolds)+len(localScaffolds))

	out = append(out, localScaffolds...)

	for _, s := range systemScaffolds {
		out = append(out, s.Root)

		for _, sub := range s.SubPackages {
			out = append(out, s.Root+"/"+sub)
		}
	}

	slices.Sort(out)
	return out, nil
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
		treelist := []printer.ListTree{}

		for _, s := range systemScaffolds {
			subs := make([]printer.ListTree, len(s.SubPackages))
			for i := range s.SubPackages {
				subs[i] = printer.ListTree{
					Text: s.SubPackages[i],
				}
			}

			treelist = append(treelist, printer.ListTree{
				Text:     s.Root,
				Children: subs,
			})
		}

		ctrl.printer.ListTree("System Scaffolds", treelist)
	}

	ctrl.printer.LineBreak()
	return nil
}
