// Package scaffold provides a simple way to scaffold a project from a template
package scaffold

import (
	"fmt"
	"io/fs"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hay-kot/scaffold/internal/core/rwfs"
)

var projectNames = [...]string{
	"{{ .Project }}",
	"{{ .ProjectSlug }}",
	"{{ .ProjectSnake }}",
	"{{ .ProjectKebab }}",
	"{{ .ProjectCamel }}",
	"templates",
}

// Project structure hold the project templates file system and configuration for
// rendering the project.
type Project struct {
	RootFS       rwfs.ReadFS
	NameTemplate string
	Name         string
	Conf         *ProjectScaffoldFile
}

func LoadProject(fileSys fs.FS) (*Project, error) {
	p := &Project{
		RootFS: fileSys,
	}

	var err error

	p.NameTemplate, err = p.validate()
	if err != nil {
		return nil, err
	}

	return p, nil
}

func (p *Project) validate() (str string, err error) {
	// Ensure there is a scaffold.yaml file
	_, err = p.RootFS.Open("scaffold.yaml")
	if err != nil {
		return "", fmt.Errorf("scaffold.yaml does not exist")
	}

	// Ensure {{ .Project }} directory exists
	for _, dir := range projectNames {
		_, err = p.RootFS.Open(dir)
		if err == nil {
			return dir, nil
		}
	}

	return "", fmt.Errorf("{{ .Project }} directory does not exist")
}

func (p *Project) AskQuestions(def map[string]string) (map[string]any, error) {
	qs := []*survey.Question{}

	if name, ok := def["Project"]; !ok {
		qs = append(qs, &survey.Question{
			Name: "Project",
			Prompt: &survey.Input{
				Message: "Project name",
			},
			Validate: survey.Required,
		})
	} else {
		p.Name = name
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

	for _, q := range scaffold.Questions {
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

	err = survey.Ask(qs, &vars)
	if err != nil {
		return nil, err
	}

	p.Name = vars["Project"].(string)

	return vars, nil
}
