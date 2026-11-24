package commands

import (
	"fmt"
	"io"
	"os"

	"charm.land/huh/v2/spinner"
	"github.com/charmbracelet/x/term"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/rs/zerolog/log"
)

// syncRCScaffolds clones the scaffoldrc scaffolds that are not cached yet. It
// never pulls; 'scaffold update' does that. A failed clone is a warning, not an
// error, so offline use keeps working with the scaffolds already cached. All
// output goes to stderr so stdout stays clean for 'scaffold list --json'.
func (ctrl *Controller) syncRCScaffolds() {
	if len(ctrl.rc.Scaffolds) == 0 {
		return
	}

	resolver := pkgs.NewResolver(ctrl.rc.Shorts, ctrl.Flags.Cache, ".", pkgs.WithProgress(io.Discard))

	var missing []string
	for _, ref := range ctrl.rc.Scaffolds {
		cached, err := resolver.Cached(ref)
		if err != nil {
			ctrl.errPrinter.Warning(fmt.Sprintf("skipping scaffold %s: %v", ref, err))
			continue
		}

		if !cached {
			missing = append(missing, ref)
		}
	}

	if len(missing) == 0 {
		return
	}

	var failures []string
	runWithStderrSpinner("Cloning scaffolds from scaffoldrc...", func() {
		for _, ref := range missing {
			if _, err := resolver.EnsureCached(ref, ctrl.rc); err != nil {
				failures = append(failures, fmt.Sprintf("could not clone scaffold %s: %v", ref, err))
			}
		}
	})

	for _, f := range failures {
		ctrl.errPrinter.Warning(f)
	}

	if len(failures) > 0 {
		ctrl.errPrinter.Warning("check the network connection and the auth entries in scaffoldrc")
	}
}

func runWithStderrSpinner(title string, work func()) {
	if !term.IsTerminal(os.Stderr.Fd()) {
		work()
		return
	}

	err := spinner.New().
		Title(title).
		WithOutput(os.Stderr).
		Action(work).
		Run()
	if err != nil {
		log.Debug().Err(err).Msg("spinner failed")
	}
}
