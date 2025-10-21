// Package commands containers all the commands for the application CLI
package commands

import (
	"bytes"
	"fmt"
	"os"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"github.com/hay-kot/scaffold/internal/printer"
	"github.com/hay-kot/scaffold/internal/styles"
	"gopkg.in/yaml.v3"
)

type Flags struct {
	ScaffoldRCPath string
	Cache          string
	ScaffoldDirs   []string
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
	ctrl.printer = printer.New(os.Stdout).WithBase(styles.Base).WithLight(styles.Light).WithWarning(styles.Warning)
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

// loadLocalScaffolds loads scaffolds from all configured scaffold directories.
// It warns when a directory doesn't exist or contains no scaffolds, but continues processing.
func (ctrl *Controller) loadLocalScaffolds() ([]string, error) {
	const Indent = " " // indent warnings to match list output
	localScaffolds := []string{}

	for _, dir := range ctrl.Flags.ScaffoldDirs {
		// Check if directory exists
		_, err := os.Stat(dir)
		if os.IsNotExist(err) {
			ctrl.printer.Warning(Indent + "Warning: scaffold directory not found: " + dir)
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("failed to check directory %s: %w", dir, err)
		}

		// List scaffolds in directory
		scaffolds, err := pkgs.ListFromFS(os.DirFS(dir))
		if err != nil {
			return nil, fmt.Errorf("failed to list scaffolds from %s: %w", dir, err)
		}

		// Warn if directory exists but has no scaffolds
		if len(scaffolds) == 0 {
			ctrl.printer.Warning(Indent + "Warning: no scaffolds found in directory: " + dir)
		}

		localScaffolds = append(localScaffolds, scaffolds...)
	}

	return localScaffolds, nil
}
