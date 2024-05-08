package commands

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/sahilm/fuzzy"
	"github.com/urfave/cli/v2"
)

func parseArgVars(args []string) (map[string]any, error) {
	vars := make(map[string]any, len(args))

	for _, v := range args {
		if !strings.Contains(v, "=") {
			return nil, fmt.Errorf("variable %s is not in the form of key=value", v)
		}

		kv := strings.Split(v, "=")
		vars[kv[0]] = kv[1]
	}

	return vars, nil
}

func (ctrl *Controller) New(ctx *cli.Context) error {
	path, err := ctrl.resolve(ctx.Args().First())
	if err != nil {
		return err
	}

	rest := ctx.Args().Tail()
	argvars, err := parseArgVars(rest)
	if err != nil {
		return err
	}

	err = ctrl.runscaffold(runconf{
		scaffolddir:  path,
		showMessages: true,
		varfunc: func(p *scaffold.Project) (map[string]any, error) {
			vars := scaffold.MergeMaps(ctrl.vars, argvars, ctrl.rc.Defaults)
			vars, err = p.AskQuestions(vars, ctrl.engine)
			if err != nil {
				return nil, err
			}

			return vars, nil
		},
		outputfs: rwfs.NewOsWFS(ctrl.Flags.OutputDir),
	})
	if err != nil {
		return err
	}

	return nil
}

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

func basicAuthAuthorizer(pkgurl, username, password string) pkgs.AuthProviderFunc {
	return func(url string) (transport.AuthMethod, bool) {
		if url != pkgurl {
			return nil, false
		}

		return &http.BasicAuth{
			Username: username,
			Password: password,
		}, true
	}
}
