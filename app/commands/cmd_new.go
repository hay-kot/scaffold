package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math/rand"
	"os"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hay-kot/scaffold/app/argparse"
	"github.com/hay-kot/scaffold/app/core/fsast"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
	"github.com/hay-kot/scaffold/internal/styles"
	"github.com/sahilm/fuzzy"
)

// DryRunOutput represents the JSON output for --dry-run
type DryRunOutput struct {
	Files    []DryRunFile `json:"files"`
	Errors   []string     `json:"errors"`
	Warnings []string     `json:"warnings"`
}

// DryRunFile represents a file that would be created
type DryRunFile struct {
	Path   string `json:"path"`
	Action string `json:"action"`
}

type FlagsNew struct {
	NoPrompt   bool
	Preset     string
	Snapshot   string
	NoClobber  bool
	ForceApply bool
	OutputDir  string
	DryRun     bool
}

// OutputFS returns a WriteFS based on the OutputDir flag.
// When DryRun is true, returns an in-memory filesystem.
func (f FlagsNew) OutputFS() rwfs.WriteFS {
	if f.OutputDir == ":memory:" || f.DryRun {
		return rwfs.NewMemoryWFS()
	}

	return rwfs.NewOsWFS(f.OutputDir)
}

func (ctrl *Controller) New(args []string, flags FlagsNew) error {
	if len(args) == 0 {
		ctrl.printer.FatalError(errors.New("missing scaffold path"))
		return ctrl.List(FlagsList{
			OutputDir: flags.OutputDir,
		})
	}

	path, err := ctrl.resolve(args[0], flags.OutputDir, flags.NoPrompt, flags.ForceApply)
	if err != nil {
		return err
	}

	if path == "" {
		return fmt.Errorf("missing scaffold path")
	}

	rest := args[1:]
	argvars, err := argparse.Parse(rest)
	if err != nil {
		return err
	}

	var varfunc func(*scaffold.Project) (map[string]any, error)
	switch {
	case flags.NoPrompt:
		varfunc = func(p *scaffold.Project) (map[string]any, error) {
			// Start with preset if specified
			var baseVars map[string]any
			if flags.Preset != "" {
				presetVars, ok := p.Conf.Presets[flags.Preset]
				if !ok {
					return nil, fmt.Errorf("preset '%s' not found", flags.Preset)
				}
				baseVars = presetVars
			} else {
				baseVars = make(map[string]any)
			}

			// Merge CLI arguments, which take precedence over presets
			vars := scaffold.MergeMaps(baseVars, argvars)

			// Ensure Project name is set
			project, ok := vars["Project"].(string)
			if !ok || project == "" {
				// Generate 4 random digits
				name := fmt.Sprintf("scaffold-test-%04d", rand.Intn(10000))
				vars["Project"] = name
				project = name
			}
			p.Name = project

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

	if flags.DryRun {
		output := DryRunOutput{
			Files:    []DryRunFile{},
			Errors:   []string{},
			Warnings: []string{},
		}

		err := fs.WalkDir(outfs, ".", func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if d.IsDir() {
				return nil
			}
			output.Files = append(output.Files, DryRunFile{
				Path:   path,
				Action: "create",
			})
			return nil
		})
		if err != nil {
			output.Errors = append(output.Errors, err.Error())
		}

		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(output)
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

			defer file.Close() //nolint:errcheck

			_, err = file.WriteString(ast.String())
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (ctrl *Controller) fuzzyFallBack(str string) ([]string, []string, error) {
	systemScaffolds, err := pkgs.ListSystem(os.DirFS(ctrl.Flags.Cache))
	if err != nil {
		return nil, nil, err
	}

	// Collect scaffolds from all configured scaffold directories
	localScaffolds, err := ctrl.loadLocalScaffolds()
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
