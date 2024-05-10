package scaffold

import (
	"io"

	"github.com/charmbracelet/huh"
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
	Presets  map[string]map[string]any `yaml:"presets"`
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
	Loop    bool      `yaml:"loop"`
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

func (p AnyPrompt) IsInputLoop() bool {
	return p.IsInput() && p.Loop
}

func (p AnyPrompt) IsMultiSelect() bool {
	return p.IsSelect() && p.Multi
}

type Askable interface {
	Ask(vars engine.Vars) error
}

type AskableFunc func(vars engine.Vars) error

func (a AskableFunc) Ask(vars engine.Vars) error {
	return a(vars)
}

func (q Question) ToHuhQuestion(defaults engine.Vars) Askable {
	def := defaults[q.Name]

	switch {
	case q.Prompt.IsMultiSelect():
		opts := make([]huh.Option[string], 0, len(*q.Prompt.Options))
		for _, option := range *q.Prompt.Options {
			opts = append(opts, huh.NewOption(option, option))
		}

		defValue := parseDefaultStrings(def, q.Prompt.Default)

		return AskableFunc(func(vars engine.Vars) error {
			ask := huh.NewMultiSelect[string]().
				Title(*q.Prompt.Message).
				Options(opts...).
				Value(&defValue)

			err := ask.Run()
			if err != nil {
				return err
			}

			vars[q.Name] = def
			return nil
		})

	case q.Prompt.IsSelect():
		opts := make([]huh.Option[string], 0, len(*q.Prompt.Options))
		for _, option := range *q.Prompt.Options {
			opts = append(opts, huh.NewOption(option, option))
		}

		defValue := parseDefaultString(def, q.Prompt.Default)

		return AskableFunc(func(vars engine.Vars) error {
			ask := huh.NewSelect[string]().
				Title(*q.Prompt.Message).
				Options(opts...).
				Value(&defValue)

			err := ask.Run()
			if err != nil {
				return err
			}

			vars[q.Name] = def
			return nil
		})
	case q.Prompt.IsConfirm():
		defValue := parseDefaultBool(def, q.Prompt.Default)
		return AskableFunc(func(vars engine.Vars) error {
			ask := huh.NewConfirm().
				Title(*q.Prompt.Confirm).
				Value(&defValue)

			err := ask.Run()
			if err != nil {
				return nil
			}

			vars[q.Name] = def
			return nil
		})
	case q.Prompt.IsInputLoop():
		var out []string

		return AskableFunc(func(vars engine.Vars) error {
			for {
				ref := ""

				ask := huh.NewInput().
					Title(q.Name).
					Value(&ref)

				err := ask.Run()
				if err != nil {
					return nil
				}

				if ref == "" {
					break
				}

				out = append(out, ref)
			}

			vars[q.Name] = out
			return nil
		})

	case q.Prompt.IsInput():
		defValue := parseDefaultString(def, q.Prompt.Default)

		return AskableFunc(func(vars engine.Vars) error {
			ask := huh.NewInput().
				Title(q.Name).
				Value(&defValue)

			err := ask.Run()
			if err != nil {
				return nil
			}

			vars[q.Name] = def
			return nil
		})
	default:
		log.Fatal().
			Str("question", q.Name).
			Msgf("Unknown prompt type")

		return nil
	}
}
