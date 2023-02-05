package scaffold

import (
	"fmt"
	"io"
	"io/fs"

	"github.com/AlecAivazis/survey/v2"
	"gopkg.in/yaml.v3"
)

var (
	projectNames = [...]string{
		"{{ .Project }}",
		"{{ .ProjectSlug }}",
		"{{ .ProjectSnake }}",
		"{{ .ProjectKebab }}",
		"{{ .ProjectCamel }}",
	}
)

// Project structure hold the project templates file system and configuration for
// rendering the project.
type Project struct {
	RootFS              fs.FS
	Tree                *TemplateNode
	ProjectNameTemplate string

	qProject string
}

func LoadProject(fileSys fs.FS) (*Project, error) {
	p := &Project{
		RootFS: fileSys,
	}

	pNameTemplate, err := p.validate()
	if err != nil {
		return nil, err
	}

	p.ProjectNameTemplate = pNameTemplate

	// Build Template Tree
	projFS, err := fs.Sub(fileSys, p.ProjectNameTemplate)
	if err != nil {
		return nil, err
	}

	tree, err := parseTemplateNodeTree(projFS, ".")
	if err != nil {
		return nil, err
	}

	p.Tree = tree

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

func (p *Project) AskQuestions() (map[string]interface{}, error) {
	qs := []*survey.Question{
		{
			Name: "Project",
			Prompt: &survey.Input{
				Message: "Project name",
			},
			Validate: survey.Required,
		},
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

	for _, q := range scaffold.Questions {
		qs = append(qs, q.ToSurveyQuestion())
	}

	vars := map[string]any{}

	err = survey.Ask(qs, &vars)
	if err != nil {
		return nil, err
	}

	p.qProject = vars["Project"].(string)

	return vars, nil
}

type Question struct {
	Name     string    `yaml:"name"`
	Prompt   AnyPrompt `yaml:"prompt"`
	Required bool      `yaml:"required"`
}

type AnyPrompt struct {
	Message *string   `yaml:"message"`
	Default string    `yaml:"default"`
	Confirm *string   `yaml:"confirm"`
	Multi   bool      `yaml:"multi"`
	Options *[]string `yaml:"options"`
}

func (p AnyPrompt) IsSelect() bool {
	return p.Message != nil && p.Options != nil
}

func (p AnyPrompt) IsConfirm() bool {
	return p.Confirm != nil
}

func (p AnyPrompt) IsInput() bool {
	return p.Message != nil
}

func (p AnyPrompt) IsMultiSelect() bool {
	return p.IsSelect() && p.Multi
}

func (q Question) ToSurveyQuestion() *survey.Question {
	out := &survey.Question{
		Name: q.Name,
	}

	switch {
	case q.Prompt.IsMultiSelect():
		out.Prompt = &survey.MultiSelect{
			Message: *q.Prompt.Message,
			Options: *q.Prompt.Options,
			Default: q.Prompt.Default,
		}
	case q.Prompt.IsSelect():
		out.Prompt = &survey.Select{
			Message: *q.Prompt.Message,
			Options: *q.Prompt.Options,
			Default: q.Prompt.Default,
		}
	case q.Prompt.IsConfirm():
		out.Prompt = &survey.Confirm{
			Message: *q.Prompt.Confirm,
			Default: q.Prompt.Default == "true",
		}
	case q.Prompt.IsInput():
		out.Prompt = &survey.Input{
			Message: *q.Prompt.Message,
			Default: q.Prompt.Default,
		}
	default:
		panic("not implemented")
	}

	return out
}

type ProjectScaffoldFile struct {
	Questions []Question `yaml:"questions"`
}

func readScaffoldFile(reader io.Reader) (*ProjectScaffoldFile, error) {
	var out ProjectScaffoldFile

	err := yaml.NewDecoder(reader).Decode(&out)
	if err != nil {
		return nil, err
	}

	return &out, nil
}
