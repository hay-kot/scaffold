// Package scaffold provides a simple way to scaffold a project from a template
package scaffold

import (
	"fmt"
	"io/fs"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
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
	scaffoldFile, err := p.RootFS.Open("scaffold.yaml")
	if err != nil {
		return nil, err
	}

	scaffold, err := readScaffoldFile(scaffoldFile)
	if err != nil {
		return nil, err
	}

	p.Conf = scaffold

	return p, nil
}

func (p *Project) validate() (str string, err error) {
	// Ensure there is a scaffold.yaml file
	_, err = p.RootFS.Open("scaffold.yaml")
	if err != nil {
		return "", fmt.Errorf("scaffold.yaml does not exist")
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

func (p *Project) AskQuestions(def map[string]string) (map[string]any, error) {
	qs := []*survey.Question{}

	projectMode := p.NameTemplate != "templates"

	if projectMode {
		name, ok := def["Project"]

		switch ok {
		case false:
			qs = append(qs, &survey.Question{
				Name: "Project",
				Prompt: &survey.Input{
					Message: "Project name",
				},
				Validate: survey.Required,
			})
		case true:
			p.Name = name
		}
	}

	for _, q := range p.Conf.Questions {
		qs = append(qs, q.ToSurveyQuestion())
	}

	// Filter out questions that have already been answered
	// TODO: Types should be checked to ensure they are what's expected in the template
	for i := 0; i < len(qs); i++ {
		if _, ok := def[qs[i].Name]; ok {
			qs = append(qs[:i], qs[i+1:]...)
			i--
		}
	}

	vars := make(map[string]any)
	// Copy default values
	for k, v := range def {
		vars[k] = v
	}

	if len(qs) == 0 {
		return vars, nil
	}

	err := survey.Ask(qs, &vars)
	if err != nil {
		return nil, err
	}

	if projectMode {
		p.Name = vars["Project"].(string)
	} else {
		p.Name = "templates"
		vars["Project"] = "templates"
	}

	// Unwrap core.OptionAnswer types into their values
	for k, v := range vars {
		switch vt := v.(type) {
		case core.OptionAnswer:
			vars[k] = vt.Value
		case []core.OptionAnswer:
			values := make([]string, len(vt))
			for i, v := range vt {
				values[i] = v.Value
			}

			vars[k] = values
		}
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
