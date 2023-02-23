// Package handlers containers all the handlers for the application CLI commands
package handlers

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
}

type Controller struct {
	engine *engine.Engine
	Flags  Flags
	// Global Flags
	rc       *scaffold.ScaffoldRC
	vars     map[string]string
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
