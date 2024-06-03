// Package commands containers all the commands for the application CLI
package commands

import (
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
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
	rc       *scaffoldrc.ScaffoldRC
	prepared bool
}

// Prepare sets up the controller to be called by the CLI, if the controller is
// not prepared it will panic
func (ctrl *Controller) Prepare(e *engine.Engine, src *scaffoldrc.ScaffoldRC) {
	ctrl.engine = e
	ctrl.rc = src
	ctrl.prepared = true
}

func (ctrl *Controller) ready() {
	if !ctrl.prepared {
		panic("controller not prepared")
	}
}
