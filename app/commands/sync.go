package commands

import (
	"fmt"
	"io"
	"os"

	"charm.land/huh/v2/spinner"
	"github.com/charmbracelet/x/term"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs/pkgurl"
)

// syncRCScaffolds clones the scaffoldrc scaffolds that are not cached yet. It
// never pulls; 'scaffold update' does that. A failed clone is a warning, not an
// error, so offline use keeps working with the scaffolds already cached. All
// output goes to stderr so stdout stays clean for 'scaffold list --json'.
//
// A spinner error, such as a Ctrl+C to skip a clone that hangs while offline,
// is also only a warning. The clone goroutine can still be running then, so
// its failures are not read; an incomplete clone is repaired on the next sync.
func (ctrl *Controller) syncRCScaffolds() {
	if len(ctrl.rc.Scaffolds) == 0 {
		return
	}

	opts := append([]pkgs.ResolverOption{pkgs.WithProgress(io.Discard)}, ctrl.syncResolverOpts...)
	resolver := pkgs.NewResolver(ctrl.rc.Shorts, ctrl.Flags.Cache, ".", opts...)

	var missing []string
	for _, ref := range ctrl.rc.Scaffolds {
		cached, err := resolver.Cached(ref)
		if err != nil {
			ctrl.errPrinter.Warning(fmt.Sprintf("skipping scaffold %s: %s", pkgurl.Redact(ref), pkgurl.Redact(err.Error())))
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
	err := runWithStderrSpinner("Cloning scaffolds from scaffoldrc...", func() {
		for _, ref := range missing {
			if err := resolver.EnsureCached(ref, ctrl.rc); err != nil {
				failures = append(failures, fmt.Sprintf("could not clone scaffold %s: %v", pkgurl.Redact(ref), err))
			}
		}
	})
	if err != nil {
		ctrl.errPrinter.Warning(fmt.Sprintf("skipped cloning scaffolds from scaffoldrc: %v", err))
		return
	}

	for _, f := range failures {
		ctrl.errPrinter.Warning(f)
	}

	if len(failures) > 0 {
		ctrl.errPrinter.Warning("check the network connection and the auth entries in scaffoldrc")
	}
}

// runWithStderrSpinner runs work behind a spinner on stderr, or directly when
// stderr is not a terminal. On a spinner error such as a Ctrl+C, work may still
// be running, so the caller must not read state that work writes.
func runWithStderrSpinner(title string, work func()) error {
	if !term.IsTerminal(os.Stderr.Fd()) {
		work()
		return nil
	}

	return spinner.New().
		Title(title).
		WithOutput(os.Stderr).
		Action(work).
		Run()
}
