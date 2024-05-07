// Package commands containers all the commands for the application CLI
package commands

import (
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold"
)

type Flags struct {
	NoClobber      bool
	Force          bool
	ScaffoldRCPath string
	Cache          string
	OutputDir      string
	ScaffoldDirs   []string
	Cwd            string
	NoInteractive  bool
}

type Controller struct {
	// Global Flags
	Flags Flags

	engine   *engine.Engine
	rc       *scaffold.ScaffoldRC
	vars     map[string]any
	prepared bool
}

// Prepare sets up the controller to be called by the CLI, if the controller is
// not prepared it will panic
func (ctrl *Controller) Prepare(e *engine.Engine, src *scaffold.ScaffoldRC) {
	ctrl.engine = e
	ctrl.rc = src
	ctrl.prepared = true
}

func (ctrl *Controller) ready() {
	if !ctrl.prepared {
		panic("controller not prepared")
	}
}
