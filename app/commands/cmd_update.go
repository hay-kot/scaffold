package commands

import (
	"context"
	"fmt"
	"os"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/internal/printer"
	"github.com/urfave/cli/v3"
)

type updateStatus struct {
	repository string
	message    string
}

func (ctrl *Controller) Update(ctx context.Context, c *cli.Command) error {
	ctrl.ready()
	ctrl.syncRCScaffolds()

	scaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return err
	}

	var failed []updateStatus
	var updated []string
	var uptodate []string
	var pinned []string

	workfn := func() {
		for _, s := range scaffolds {
			status, err := pkgs.Update(ctrl.Flags.Cache, s.Root, ctrl.rc)
			if err != nil {
				failed = append(failed, updateStatus{
					repository: s.Root,
					message:    err.Error(),
				})

				continue
			}

			switch status {
			case pkgs.UpdateStatusPinned:
				pinned = append(pinned, s.Root)
			case pkgs.UpdateStatusUpdated:
				updated = append(updated, s.Root)
			default:
				uptodate = append(uptodate, s.Root)
			}
		}
	}

	runWithStderrSpinner("Updating Scaffolds...", workfn)

	if len(uptodate) > 0 {
		ctrl.printer.Title(fmt.Sprintf("%d scaffolds up to date", len(uptodate)))
	}

	if len(pinned) > 0 {
		ctrl.printer.List("Pinned, not pulled", pinned)
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
		items := make([]printer.StatusListItem, 0, len(failed))
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
