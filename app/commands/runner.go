package commands

import (
	"fmt"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
)

type runconf struct {
	// os path to the scaffold directory.
	scaffolddir string
	// showMessages is a flag to show pre/post messages.
	showMessages bool
	// varfunc is a function that returns a map of variables that is provided
	// to the template engine.
	varfunc func(*scaffold.Project) (map[string]any, error)
	// outputdir is the output directory or filesystem.
	outputfs rwfs.WriteFS
}

// runscaffold runs the scaffold. This method exists outside of the `new` receiver function
// so that we can allow the `test` and `new` commands to share as much of the same code
// as possible.
func (ctrl *Controller) runscaffold(cfg runconf) error {
	pfs := os.DirFS(cfg.scaffolddir)
	p, err := scaffold.LoadProject(pfs, scaffold.Options{
		NoClobber: ctrl.Flags.NoClobber,
	})
	if err != nil {
		return err
	}

	if cfg.showMessages && p.Conf.Messages.Pre != "" {
		out, err := glamour.RenderWithEnvironmentConfig(p.Conf.Messages.Pre)
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	vars, err := cfg.varfunc(p)
	if err != nil {
		return err
	}

	args := &scaffold.RWFSArgs{
		Project: p,
		ReadFS:  pfs,
		WriteFS: cfg.outputfs,
	}

	vars, err = scaffold.BuildVars(ctrl.engine, args.Project, vars)
	if err != nil {
		return err
	}

	err = scaffold.RenderRWFS(ctrl.engine, args, vars)
	if err != nil {
		return err
	}

	err = ctrl.RunHook(pfs, "post_scaffold", cfg.outputfs, vars)
	if err != nil {
		return err
	}

	if cfg.showMessages && p.Conf.Messages.Post != "" {
		rendered, err := ctrl.engine.TmplString(p.Conf.Messages.Post, vars)
		if err != nil {
			return err
		}

		out, err := glamour.RenderWithEnvironmentConfig(rendered)
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	return nil
}
