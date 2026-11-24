package commands

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/signal"
	"sync/atomic"

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
// Ctrl+C, for example to skip a clone that hangs while offline, is also only a
// warning. It stops the clone in progress and removes its directory, so the
// command continues with a cache that holds no partial clone.
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
	err := runWithStderrSpinner("Cloning scaffolds from scaffoldrc...", func(ctx context.Context) {
		for _, ref := range missing {
			if ctx.Err() != nil {
				return
			}

			err := resolver.EnsureCached(ctx, ref, ctrl.rc)
			if err != nil && ctx.Err() == nil {
				failures = append(failures, fmt.Sprintf("could not clone scaffold %s: %v", pkgurl.Redact(ref), err))
			}
		}
	})

	for _, f := range failures {
		ctrl.errPrinter.Warning(f)
	}

	if err != nil {
		ctrl.errPrinter.Warning("interrupted, skipped cloning the remaining scaffolds from scaffoldrc")
		return
	}

	if len(failures) > 0 {
		ctrl.errPrinter.Warning("check the network connection and the auth entries in scaffoldrc")
	}
}

var errInterrupted = errors.New("interrupted")

// runWithStderrSpinner runs work behind a spinner on stderr, or directly when
// stderr is not a terminal. Ctrl+C (a key press under the spinner, SIGINT
// without it) cancels the ctx that work gets. runWithStderrSpinner then waits
// for work to return, so work can clean up and its results are safe to read,
// and returns an error.
func runWithStderrSpinner(title string, work func(ctx context.Context)) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var started atomic.Bool
	done := make(chan struct{})
	run := func() {
		started.Store(true)
		defer close(done)
		work(ctx)
	}

	if !term.IsTerminal(os.Stderr.Fd()) {
		run()
		if ctx.Err() != nil {
			return errInterrupted
		}
		return nil
	}

	err := spinner.New().
		Title(title).
		WithOutput(os.Stderr).
		Action(run).
		Run()
	if err != nil {
		cancel()
		// The spinner can fail before it starts the action.
		if started.Load() {
			<-done
		}
		return err
	}

	return nil
}
