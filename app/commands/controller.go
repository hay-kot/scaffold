// Package commands containers all the commands for the application CLI
package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rule"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
)

type Flags struct {
	NoClobber      bool
	Force          bool
	ScaffoldRCPath string
	Cache          string
	OutputDir      string
	ScaffoldDirs   []string
}

// OutputFS returns a WriteFS based on the OutputDir flag
func (f Flags) OutputFS() rwfs.WriteFS {
	if f.OutputDir == ":memory:" {
		return rwfs.NewMemoryWFS()
	}

	return rwfs.NewOsWFS(f.OutputDir)
}

type Controller struct {
	// Flags contains the CLI flags
	// that are from the root command
	Flags Flags

	engine   *engine.Engine
	rc       *scaffold.ScaffoldRC
	runHooks rule.Rule
	prepared bool
}

// Prepare sets up the controller to be called by the CLI, if the controller is
// not prepared it will panic
func (ctrl *Controller) Prepare(e *engine.Engine, src *scaffold.ScaffoldRC) {
	ctrl.engine = e
	ctrl.rc = src
	ctrl.prepared = true
	ctrl.runHooks = src.RunHooks
}

func (ctrl *Controller) RunHook(rfs rwfs.ReadFS, name string, wfs rwfs.WriteFS, vars any, args ...string) error {
	src, err := fs.ReadFile(rfs, filepath.Join("hooks", name))
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to open %q hook file: %w", name, err)
		}
		return nil
	}

	rendered, err := ctrl.engine.TmplString(string(src), vars)
	if err != nil {
		return err
	}

	if !mayRunHook(ctrl.runHooks, name, rendered) {
		return nil
	}

	err = wfs.RunHook(name, []byte(rendered), args)
	if err != nil && !errors.Is(err, rwfs.ErrHooksNotSupported) {
		return err
	}

	return nil
}

func (ctrl *Controller) ready() {
	if !ctrl.prepared {
		panic("controller not prepared")
	}
}

func mayRunHook(hookRule rule.Rule, name string, rendered string) bool {
	for {
		switch hookRule {
		case rule.Unset:
			return true
		case rule.Yes:
			return true
		case rule.No:
			return false
		case rule.Prompt:
		}

		err := huh.Run(huh.NewSelect[rule.Rule]().
			Title(fmt.Sprintf("scaffold defines a %s hook", name)).
			Options(
				huh.NewOption("run", rule.Yes),
				huh.NewOption("skip", rule.No),
				huh.NewOption("review", rule.Prompt)).
			Value(&hookRule))

		if err != nil {
			fmt.Fprint(os.Stderr, err)
			return false
		}

		if hookRule != rule.Prompt {
			continue
		}

		fmt.Printf("\n%s\n", rendered)
	}
}
