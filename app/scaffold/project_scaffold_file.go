package scaffold

import (
	"io"

	"github.com/AlecAivazis/survey/v2"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type ProjectScaffoldFile struct {
	Skip      []string          `yaml:"skip"`
	Questions []Question        `yaml:"questions"`
	Rewrites  []Rewrite         `yaml:"rewrites"`
	Computed  map[string]string `yaml:"computed"`
	Messages  struct {
		Pre  string `yaml:"pre"`
		Post string `yaml:"post"`
	} `yaml:"messages"`
	Inject   []Injectable              `yaml:"inject"`
	Features []Feature                 `yaml:"features"`
	Tests    map[string]map[string]any `yaml:"tests"`
}

type Rewrite struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Injectable struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	At       string `yaml:"at"`
	Template string `yaml:"template"`
}

type Feature struct {
	Value string   `yaml:"value"`
	Globs []string `yaml:"globs"`
}

func ReadScaffoldFile(reader io.Reader) (*ProjectScaffoldFile, error) {
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
	When     string    `yaml:"when"`
	Required bool      `yaml:"required"`
}

type AnyPrompt struct {
	Message *string   `yaml:"message"`
	Default any       `yaml:"default"`
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

// parseDefaults parses the default values in order of priority:
// where the first argument has the highest priority. As soon as a
// type match is found, the value is returned.
//
// If no match is found, the default value is returned.
func parseDefaults[T any](v ...any) T {
	for _, val := range v {
		log.Debug().Type("val", val).Msg("parseDefaults")
		val, ok := val.(T)
		if ok {
			return val
		}
	}

	var out T
	return out
}

func (q Question) ToSurveyQuestion(defaults engine.Vars) *survey.Question {
	out := &survey.Question{
		Name: q.Name,
	}

	def := defaults[q.Name]

	switch {
	case q.Prompt.IsMultiSelect():
		out.Prompt = &survey.MultiSelect{
			Message: *q.Prompt.Message,
			Options: *q.Prompt.Options,
			Default: parseDefaults[[]any](def, q.Prompt.Default),
		}
	case q.Prompt.IsSelect():
		out.Prompt = &survey.Select{
			Message: *q.Prompt.Message,
			Options: *q.Prompt.Options,
			Default: parseDefaults[string](def, q.Prompt.Default),
		}
	case q.Prompt.IsConfirm():
		out.Prompt = &survey.Confirm{
			Message: *q.Prompt.Confirm,
			Default: parseDefaults[bool](def, q.Prompt.Default),
		}
	case q.Prompt.IsInput():
		out.Prompt = &survey.Input{
			Message: *q.Prompt.Message,
			Default: parseDefaults[string](def, q.Prompt.Default),
		}
	default:
		log.Fatal().
			Str("question", q.Name).
			Msgf("Unknown prompt type")
	}

	if q.Required {
		out.Validate = survey.Required
	}

	return out
}
