package commands

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
)

type runconf struct {
	// os path to the scaffold directory.
	scaffolddir string
	// noPrompt is a flag to show pre/post messages.
	noPrompt bool
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
	scaffoldFS := os.DirFS(cfg.scaffolddir)
	p, err := scaffold.LoadProject(scaffoldFS, scaffold.Options{
		NoClobber: ctrl.Flags.NoClobber,
	})
	if err != nil {
		return err
	}

	if !cfg.noPrompt && p.Conf.Messages.Pre != "" {
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
		ReadFS:  scaffoldFS,
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

	if ctrl.rc.Settings.RunHooks != scaffoldrc.RunHooksNever {
		err = ctrl.runHook(scaffoldFS, cfg.outputfs, scaffold.PostScaffoldScripts, vars, cfg.noPrompt)
		if err != nil {
			return err
		}
	}

	if !cfg.noPrompt && p.Conf.Messages.Post != "" {
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

func (ctrl *Controller) runHook(
	rfs rwfs.ReadFS,
	wfs rwfs.WriteFS,
	hookPrefix string,
	vars any,
	noPrompt bool,
) error {
	sources, err := fs.ReadDir(rfs, scaffold.HooksDir)
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("failed to open %q hook file: %w", hookPrefix, err)
		}
		return nil
	}

	// find first glob match
	var hookContents []byte
	for _, source := range sources {
		if source.IsDir() {
			continue
		}

		if strings.HasPrefix(source.Name(), hookPrefix) {
			path := filepath.Join(scaffold.HooksDir, source.Name())

			hookContents, err = fs.ReadFile(rfs, path)
			if err != nil {
				return err
			}
		}
	}

	if len(hookContents) == 0 {
		return nil
	}

	rendered, err := ctrl.engine.TmplString(string(hookContents), vars)
	if err != nil {
		return err
	}

	if !shouldRunHooks(ctrl.rc.Settings.RunHooks, noPrompt, hookPrefix, rendered) {
		return nil
	}

	err = wfs.RunHook(hookPrefix, []byte(rendered), nil)
	if err != nil && !errors.Is(err, rwfs.ErrHooksNotSupported) {
		return err
	}

	return nil
}

// shouldRunHooks will resolve the users RunHooks preference and either return the preference
// or prompt the user for their choice when the preference is RunHooksPrompt
func shouldRunHooks(runPreference scaffoldrc.RunHooksOption, noPrompt bool, name string, rendered string) bool {
	for {
		switch runPreference {
		case scaffoldrc.RunHooksAlways:
			return true
		case scaffoldrc.RunHooksNever:
			return false
		case scaffoldrc.RunHooksPrompt:
			if noPrompt {
				return false
			}

			err := huh.Run(huh.NewSelect[scaffoldrc.RunHooksOption]().
				Title(fmt.Sprintf("scaffold defines a %s hook", name)).
				Options(
					huh.NewOption("run", scaffoldrc.RunHooksAlways),
					huh.NewOption("skip", scaffoldrc.RunHooksNever),
					huh.NewOption("review", scaffoldrc.RunHooksPrompt)).
				Value(&runPreference))
			if err != nil {
				fmt.Fprint(os.Stderr, err)
				return false
			}

			if runPreference == scaffoldrc.RunHooksPrompt {
				fmt.Printf("\n%s\n", rendered)
			}
		default:
			return false
		}
	}
}
