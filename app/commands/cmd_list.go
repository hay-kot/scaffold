package commands

import (
	"encoding/json"
	"os"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/internal/printer"
)

type FlagsList struct {
	OutputDir string
	JSON      bool
}

// ListOutput is the JSON output format for the list command.
type ListOutput struct {
	Local  []string           `json:"local"`
	System []ListSystemOutput `json:"system"`
}

// ListSystemOutput represents a system scaffold with its subpackages.
type ListSystemOutput struct {
	Root        string   `json:"root"`
	SubPackages []string `json:"subpackages,omitempty"`
}

func (ctrl *Controller) List(flags FlagsList) error {
	systemScaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return err
	}

	// Collect scaffolds from all configured scaffold directories
	localScaffolds, err := ctrl.loadLocalScaffolds()
	if err != nil {
		return err
	}

	if flags.JSON {
		return ctrl.listJSON(localScaffolds, systemScaffolds)
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

func (ctrl *Controller) listJSON(localScaffolds []string, systemScaffolds []pkgs.PackageList) error {
	output := ListOutput{
		Local:  localScaffolds,
		System: make([]ListSystemOutput, len(systemScaffolds)),
	}

	if output.Local == nil {
		output.Local = []string{}
	}

	for i, s := range systemScaffolds {
		output.System[i] = ListSystemOutput{
			Root:        s.Root,
			SubPackages: s.SubPackages,
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}
