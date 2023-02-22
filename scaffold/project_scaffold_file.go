package scaffold

import (
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Rewrite struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type ProjectScaffoldFile struct {
	Skip      []string          `yaml:"skip"`
	Questions []Question        `yaml:"questions"`
	Rewrites  []Rewrite         `yaml:"rewrites"`
	Computed  map[string]string `yaml:"computed"`
}

func readScaffoldFile(reader io.Reader) (*ProjectScaffoldFile, error) {
	var out ProjectScaffoldFile

	err := yaml.NewDecoder(reader).Decode(&out)
	if err != nil {
		return nil, err
	}

	return &out, nil
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
		log.Fatal().
			Str("question", q.Name).
			Msgf("Unknown prompt type")
	}

	return out
}
