package commands

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hay-kot/scaffold/app/core/fsast"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/internal/styles"
	"github.com/sahilm/fuzzy"
)

type FlagsNew struct {
	NoPrompt   bool
	Preset     string
	Snapshot   string
	NoClobber  bool
	ForceApply bool
	OutputDir  string
	RunHooks   bool
}

// OutputFS returns a WriteFS based on the OutputDir flag
func (f FlagsNew) OutputFS() rwfs.WriteFS {
	if f.OutputDir == ":memory:" {
		return rwfs.NewMemoryWFS()
	}

	return rwfs.NewOsWFS(f.OutputDir)
}

func (ctrl *Controller) New(args []string, flags FlagsNew) error {
	if len(args) == 0 {
		return fmt.Errorf("missing scaffold name")
	}

	path, err := ctrl.resolve(args[0], flags.OutputDir, flags.NoPrompt, flags.ForceApply)
	if err != nil {
		return err
	}

	if path == "" {
		return fmt.Errorf("missing scaffold path")
	}

	rest := args[1:]
	argvars, err := parseArgVars(rest)
	if err != nil {
		return err
	}

	var varfunc func(*scaffold.Project) (map[string]any, error)
	switch {
	case flags.NoPrompt:
		varfunc = func(p *scaffold.Project) (map[string]any, error) {
			caseVars, ok := p.Conf.Presets[flags.Preset]
			if !ok {
				return nil, fmt.Errorf("preset '%s' not found", flags.Preset)
			}

			project, ok := caseVars["Project"].(string)
			if !ok || project == "" {
				// Generate 4 random digits
				name := fmt.Sprintf("scaffold-test-%04d", rand.Intn(10000))
				caseVars["Project"] = name
				project = name
			}
			p.Name = project

			// Test cases do not use rc.Defaults
			vars := scaffold.MergeMaps(caseVars, argvars)
			return vars, nil
		}

	default:
		varfunc = func(p *scaffold.Project) (map[string]any, error) {
			vars := scaffold.MergeMaps(argvars, ctrl.rc.Defaults)
			vars, err = p.AskQuestions(vars, ctrl.engine, styles.Theme(ctrl.rc.Settings.Theme))
			if err != nil {
				return nil, err
			}

			return vars, nil
		}
	}

	outfs := flags.OutputFS()

	err = ctrl.runscaffold(runconf{
		scaffolddir: path,
		noPrompt:    flags.NoPrompt,
		varfunc:     varfunc,
		outputfs:    outfs,
		options: scaffold.Options{
			NoClobber: flags.NoClobber,
		},
	})
	if err != nil {
		return err
	}

	if flags.Snapshot != "" {
		ast, err := fsast.New(outfs)
		if err != nil {
			return err
		}

		if flags.Snapshot == "stdout" {
			fmt.Println(ast.String())
		} else {
			file, err := os.Create(flags.Snapshot)
			if err != nil {
				return err
			}

			_ = file.Close()

			_, err = file.WriteString(ast.String())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ctrl *Controller) fuzzyFallBack(str, outputdir string) ([]string, []string, error) {
	systemScaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return nil, nil, err
	}

	localScaffolds, err := pkgs.ListLocal(os.DirFS(outputdir))
	if err != nil {
		return nil, nil, err
	}

	systemScaffoldsStrings := make([]string, 0, len(systemScaffolds))
	for _, s := range systemScaffolds {
		if len(s.SubPackages) > 0 {
			for _, sub := range s.SubPackages {
				systemScaffoldsStrings = append(systemScaffoldsStrings, fmt.Sprintf("%s#%s", s.Root, sub))
			}
		} else {
			systemScaffoldsStrings = append(systemScaffoldsStrings, s.Root)
		}
	}

	systemMatches := fuzzy.Find(str, systemScaffoldsStrings)
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
