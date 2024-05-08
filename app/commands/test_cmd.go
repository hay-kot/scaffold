package commands

import (
	"fmt"
	"math/rand"

	"github.com/hay-kot/scaffold/app/core/fsast"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
)

type FlagsTest struct {
	Case  string `json:"case"`
	MemFS bool   `json:"memfs"`
	Ast   bool   `json:"ast"`
}

func (ctrl *Controller) Test(args []string, flags FlagsTest) error {
	if len(args) == 0 {
		return fmt.Errorf("missing scaffold path")
	}

	path, err := ctrl.resolve(args[0])
	if err != nil {
		return err
	}

	rest := args[1:]
	argvars, err := parseArgVars(rest)
	if err != nil {
		return err
	}

	var outfs rwfs.WriteFS
	if flags.MemFS {
		outfs = rwfs.NewMemoryWFS()
	} else {
		outfs = rwfs.NewOsWFS(ctrl.Flags.OutputDir)
	}

	err = ctrl.runscaffold(runconf{
		scaffolddir:  path,
		showMessages: false,
		varfunc: func(p *scaffold.Project) (map[string]any, error) {
			caseVars, ok := p.Conf.Tests[flags.Case]
			if !ok {
				return nil, fmt.Errorf("case %s not found", flags.Case)
			}

			project, ok := caseVars["Project"].(string)
			if !ok || project == "" {
				// Generate 4 random digits
				name := fmt.Sprintf("scaffold-test-%04d", rand.Intn(10000))
				p.Name = name
				caseVars["Project"] = name
			}

			// Test cases do not use rc.Defaults
			vars := scaffold.MergeMaps(ctrl.vars, caseVars, argvars)
			return vars, nil
		},
		outputfs: outfs,
	})
	if err != nil {
		return err
	}

	if flags.Ast {
		ast, err := fsast.New(outfs)
		if err != nil {
			return err
		}

		fmt.Println(ast.String())
	}

	return nil
}
