package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/huh/spinner"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/internal/printer"
	"github.com/urfave/cli/v2"
)

type updateStatus struct {
	repository string
	message    string
}

func (ctrl *Controller) Update(ctx *cli.Context) error {
	ctrl.ready()

	scaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return err
	}

	var failed []updateStatus
	var updated []string
	var uptodate []string

	workfn := func() {
		for _, s := range scaffolds {
			isUpdated, err := pkgs.Update(ctrl.Flags.Cache, s.Root)
			if err != nil {
				failed = append(failed, updateStatus{
					repository: s.Root,
					message:    err.Error(),
				})

				continue
			}

			if !isUpdated {
				uptodate = append(uptodate, s.Root)

				continue
			}

			updated = append(updated, s.Root)
		}
	}

	err = spinner.New().
		Title("Updating Scaffolds...").
		Action(workfn).
		Run()
	if err != nil {
		return err
	}

	if len(uptodate) > 0 {
		ctrl.printer.Title(fmt.Sprintf("%d scaffolds up to date", len(uptodate)))
	}

	if len(updated) > 0 {
		items := make([]printer.StatusListItem, 0, len(updated))
		for _, s := range updated {
			items = append(items, printer.StatusListItem{
				Ok:     true,
				Status: s,
			})
		}

		ctrl.printer.StatusList("Updated", items)
	}

	if len(failed) > 0 {
		items := make([]printer.StatusListItem, 0, len(updated))
		for _, s := range failed {
			items = append(items, printer.StatusListItem{
				Ok:     false,
				Status: s.repository + ": " + s.message,
			})
		}

		ctrl.printer.StatusList("Failed", items)
	}

	return nil
}
