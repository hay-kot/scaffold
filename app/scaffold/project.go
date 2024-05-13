// Package scaffold provides a simple way to scaffold a project from a template
package scaffold

import (
	"fmt"
	"io/fs"
	"maps"
	"strconv"

	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var projectNames = [...]string{
	"{{ .Project }}",
	"{{ .ProjectSlug }}",
	"{{ .ProjectSnake }}",
	"{{ .ProjectKebab }}",
	"{{ .ProjectCamel }}",
}

var templateNames = [...]string{
	"templates",
}

// Project structure hold the project templates file system and configuration for
// rendering the project.
type Project struct {
	RootFS       rwfs.ReadFS
	NameTemplate string
	Name         string
	Conf         *ProjectScaffoldFile
	Options      Options
}

func readFirst(fsys fs.FS, names ...string) (fs.File, error) {
	for _, name := range names {
		file, err := fsys.Open(name)
		if err == nil {
			return file, nil
		}
	}

	return nil, fmt.Errorf("file not found")
}

func LoadProject(fileSys fs.FS, opts Options) (*Project, error) {
	p := &Project{
		RootFS:  fileSys,
		Options: opts,
	}

	var err error

	p.NameTemplate, err = p.validate()
	if err != nil {
		return nil, err
	}

	// Read scaffold.yaml
	scaffoldFile, err := readFirst(p.RootFS, "scaffold.yaml", "scaffold.yml")
	if err != nil {
		return nil, err
	}

	scaffold, err := ReadScaffoldFile(scaffoldFile)
	if err != nil {
		return nil, err
	}

	p.Conf = scaffold

	return p, nil
}

func (p *Project) validate() (str string, err error) {
	// Ensure there is a scaffold.yaml file
	_, err = readFirst(p.RootFS, "scaffold.yaml", "scaffold.yml")
	if err != nil {
		return "", fmt.Errorf("scaffold.{yml,yaml} does not exist")
	}

	// Ensure required directories exist
	for _, dir := range append(projectNames[:], templateNames[:]...) {
		_, err = p.RootFS.Open(dir)
		if err == nil {
			return dir, nil
		}
	}

	return "", fmt.Errorf("{{ .Project }} directory does not exist")
}

func (p *Project) AskQuestions(def map[string]any, e *engine.Engine) (map[string]any, error) {
	projectMode := p.NameTemplate != "templates"

	if projectMode {
		name, ok := def["Project"]

		switch ok {
		case false:
			msg := "Project Name"
			decs := "The name your project will be generated with"

			pre := []Question{
				{
					Name: "Project",
					Prompt: AnyPrompt{
						Message:    &msg,
						Desciption: &decs,
					},
					Required: true,
				},
			}

			p.Conf.Questions = append(pre, p.Conf.Questions...)
		case true:
			nameStr, ok := name.(string)
			if !ok {
				return nil, fmt.Errorf("Project name must be a string")
			}

			p.Name = nameStr
		}
	}

	vars := maps.Clone(def)

	qgroups := QuestionGroupBy(p.Conf.Questions)
	askables := []*Askable{}
	patchvars := func() error {
		for _, askable := range askables {
			err := askable.Hook(vars)
			if err != nil {
				return err
			}
		}

		return nil
	}

	var form *huh.Form
	formgroups := []*huh.Group{}

	for _, qgroup := range qgroups {
		fields := []huh.Field{}

		for _, q := range qgroup {
			question := q.ToAskable(vars[q.Name])
			fields = append(fields, question.Field)
			askables = append(askables, question)
		}

		group := huh.NewGroup(fields...)

		firstq := qgroup[0]
		if firstq.When != "" {
			group.WithHideFunc(func() bool {
				if form == nil {
					return false
				}

				// extract existing properties
				_ = patchvars()

				first := qgroup[0]

				// we check the first question in the group to see if it has a when
				// and if so, we evaluate it and skip the group if it's false
				if first.When != "" {
					result, err := e.TmplString(first.When, vars)
					if err != nil {
						return true
					}

					resultBool, _ := strconv.ParseBool(result)
					if !resultBool {
						return true
					}
				}
				return false
			})
		}

		formgroups = append(formgroups, group)
	}

	form = huh.NewForm(formgroups...)

	err := form.Run()
	if err != nil {
		return nil, err
	}

	// Ensure properts are set on vars
	err = patchvars()
	if err != nil {
		return nil, err
	}

	// Grab the project name from the vars/answers to ensure that
	// it's set.
	if projectMode {
		p.Name = vars["Project"].(string)
	} else {
		p.Name = "templates"
		vars["Project"] = "templates"
	}

	if log.Logger.GetLevel() == zerolog.DebugLevel {
		for k, v := range vars {
			log.Debug().
				Str("key", k).
				Type("type", v).
				Interface("value", v).
				Msg("question")
		}
	}

	for _, askable := range askables {
		fmt.Print(askable.String())
	}

	fmt.Print("\n")

	return vars, nil
}
