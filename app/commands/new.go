package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/rs/zerolog/log"
	"github.com/sahilm/fuzzy"
	"github.com/urfave/cli/v2"
)

func (ctrl *Controller) fuzzyFallBack(str string) ([]string, []string, error) {
	systemScaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return nil, nil, err
	}

	localScaffolds, err := pkgs.ListLocal(os.DirFS(ctrl.Flags.OutputDir))
	if err != nil {
		return nil, nil, err
	}

	systemMatches := fuzzy.Find(str, systemScaffolds)
	systemMatchesOutput := make([]string, len(systemMatches))
	for i, match := range systemMatches {
		systemMatchesOutput[i] = match.Str
	}

	localMatches := fuzzy.Find(str, localScaffolds)
	localMatchesOutput := make([]string, len(localMatches))
	for i, match := range localMatches {
		localMatchesOutput[i] = match.Str
	}

	return systemMatchesOutput, localMatchesOutput, nil
}

var (
	bold     = lipgloss.NewStyle().Bold(true)
	colorRed = lipgloss.NewStyle().Foreground(lipgloss.Color("#dc2626"))
)

func didYouMeanPrompt(given, suggestion string) bool {
	bldr := strings.Builder{}

	// Couldn't find a scaffold named:
	//   'foo'
	//
	// Did you mean:
	//   'bar'?
	//
	// [y/n]:

	bldr.WriteString("\n ")
	bldr.WriteString(bold.Render(colorRed.Render("could not find a scaffold named")))
	bldr.WriteString("\n    ")
	bldr.WriteString(given)
	bldr.WriteString("\n\n")
	bldr.WriteString(" ")
	bldr.WriteString(bold.Render("did you mean"))
	bldr.WriteString("\n    ")
	bldr.WriteString(suggestion)
	bldr.WriteString("?\n\n ")
	bldr.WriteString("[y/n]: ")

	out := bldr.String()

	var resp string

	fmt.Print(out)
	fmt.Scanln(&resp)

	return resp == "y"
}

func (ctrl *Controller) Project(ctx *cli.Context) error {
	argPath := ctx.Args().First()
	if argPath == "" {
		return fmt.Errorf("path is required")
	}

	// Status() call for go-git is too slow to be used here
	// https://github.com/go-git/go-git/issues/181
	if !ctrl.Flags.Force {
		ok := checkWorkingTree(ctrl.Flags.OutputDir)
		if !ok {
			log.Warn().Msg("working tree is dirty, use --force to apply changes")
			return nil
		}
	}

	resolver := pkgs.NewResolver(ctrl.rc.Shorts, ctrl.Flags.Cache, ".")

	if v, ok := ctrl.rc.Aliases[argPath]; ok {
		argPath = v
	}

	path, err := resolver.Resolve(argPath, ctrl.Flags.ScaffoldDirs)
	if err != nil {
		orgErr := err
		systemMatches, localMatches, err := ctrl.fuzzyFallBack(argPath)
		if err != nil {
			return err
		}

		var first string
		var isSystemMatch bool
		if len(systemMatches) > 0 {
			first = systemMatches[0]
			isSystemMatch = true
		}

		if len(localMatches) > 0 {
			first = localMatches[0]
		}

		if first != "" {
			useMatch := didYouMeanPrompt(argPath, first)

			if useMatch {
				if isSystemMatch {
					// prepend https:// so it resolves to the correct path
					first = "https://" + first
				}

				resolved, err := resolver.Resolve(first, ctrl.Flags.ScaffoldDirs)
				if err != nil {
					return err
				}

				path = resolved
			}
		}

		if path == "" {
			return fmt.Errorf("failed to resolve path: %w", orgErr)
		}
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
