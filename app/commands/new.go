package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
)

func (ctrl *Controller) Project(ctx *cli.Context) error {
	argPath := ctx.Args().First()
	if argPath == "" {
		return fmt.Errorf("path is required")
	}

	if !ctrl.Flags.Force {
		ok := checkWorkingTree(ctrl.Flags.Cache)
		if !ok {
			log.Warn().Msg("working tree is dirty, use --force to apply changes")
			return nil
		}
	}

	resolver := pkgs.NewResolver(
		ctrl.rc.Shorts,
		ctrl.Flags.Cache,
		".",
	)

	if v, ok := ctrl.rc.Aliases[argPath]; ok {
		argPath = v
	}

	path, err := resolver.Resolve(argPath, ctrl.Flags.ScaffoldDirs)
	if err != nil {
		return fmt.Errorf("failed to resolve path: %w", err)
	}

	rest := ctx.Args().Tail()

	ctrl.vars = make(map[string]string, len(rest))
	for _, v := range rest {
		kv := strings.Split(v, "=")
		ctrl.vars[kv[0]] = kv[1]
	}

	pfs := os.DirFS(path)

	p, err := scaffold.LoadProject(pfs, scaffold.Options{
		NoClobber: ctrl.Flags.NoClobber,
	})
	if err != nil {
		return err
	}

	defaults := scaffold.MergeMaps(ctrl.vars, ctrl.rc.Defaults)

	if p.Conf.Messages.Pre != "" {
		out, err := glamour.RenderWithEnvironmentConfig(p.Conf.Messages.Pre)
		if err != nil {
			return err
		}

		fmt.Println(out)
	}

	vars, err := p.AskQuestions(defaults, ctrl.engine)
	if err != nil {
		return err
	}

	args := &scaffold.RWFSArgs{
		Project: p,
		ReadFS:  pfs,
		WriteFS: rwfs.NewOsWFS(ctrl.Flags.OutputDir),
	}

	vars, err = scaffold.BuildVars(ctrl.engine, args, vars)
	if err != nil {
		return err
	}

	err = scaffold.RenderRWFS(ctrl.engine, args, vars)

	if err != nil {
		return err
	}

	if p.Conf.Messages.Post != "" {
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
