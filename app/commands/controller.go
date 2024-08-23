// Package commands containers all the commands for the application CLI
package commands

import (
	"bytes"
	"os"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"github.com/hay-kot/scaffold/internal/printer"
	"github.com/hay-kot/scaffold/internal/styles"
	"gopkg.in/yaml.v3"
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
	Flags   Flags
	Version string

	prepared bool
	engine   *engine.Engine
	rc       *scaffoldrc.ScaffoldRC
	printer  *printer.Printer
}

// Prepare sets up the controller to be called by the CLI, if the controller is
// not prepared it will panic
func (ctrl *Controller) Prepare(e *engine.Engine, src *scaffoldrc.ScaffoldRC) {
	ctrl.engine = e
	ctrl.rc = src
	ctrl.prepared = true
	ctrl.printer = printer.New(os.Stdout).WithBase(styles.Base).WithLight(styles.Light)
}

func (ctrl *Controller) ready() {
	if !ctrl.prepared {
		panic("controller not prepared")
	}
}

func (ctrl *Controller) RuntimeConfigYAML() (string, error) {
	buff := bytes.NewBuffer(nil)
	encoder := yaml.NewEncoder(buff)
	encoder.SetIndent(4)
	err := encoder.Encode(ctrl.rc)
	if err != nil {
		return "", err
	}

	return buff.String(), nil
}
