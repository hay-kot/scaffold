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
			msg := "Project name"
			decs := "The project name for the new project."

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

	for _, q := range p.Conf.Questions {
		if q.When != "" {
			result, err := e.TmplString(q.When, vars)
			if err != nil {
				return nil, err
			}

			resultBool, _ := strconv.ParseBool(result)
			if !resultBool {
				continue
			}
		}

		question := q.ToAskable2(vars)

		form := huh.NewForm(
			huh.NewGroup(
				question.Field,
			),
		)

		err := form.Run()
		if err != nil {
			return nil, err
		}

		err = question.Hook(vars)
		if err != nil {
			return nil, err
		}
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

	return vars, nil
}
