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

func toHuhOptions(opts *[]string) []huh.Option[string] {
	out := make([]huh.Option[string], len(*opts))
	for i, opt := range *opts {
		out[i] = huh.NewOption(opt, opt)
	}
	return out
}

func (q Question) ToAskable(defaults engine.Vars) Askable {
	def := defaults[q.Name]

	switch {
	case q.Prompt.IsMultiSelect():
		defValue := parseDefaultStrings(def, q.Prompt.Default)

		return HuhToAskable[[]string](q.Name, huh.NewMultiSelect[string]().
			Title(*q.Prompt.Message).
			Options(toHuhOptions(q.Prompt.Options)...).
			Value(&defValue))

	case q.Prompt.IsSelect():
		defValue := parseDefaultString(def, q.Prompt.Default)

		return HuhToAskable[string](q.Name, huh.NewSelect[string]().
			Title(*q.Prompt.Message).
			Options(toHuhOptions(q.Prompt.Options)...).
			Value(&defValue))
	case q.Prompt.IsConfirm():
		defValue := parseDefaultBool(def, q.Prompt.Default)
		return HuhToAskable[bool](q.Name, huh.NewConfirm().
			Title(*q.Prompt.Message).
			Value(&defValue))
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

		return HuhToAskable[string](q.Name, huh.NewInput().
			Title(*q.Prompt.Message).
			Value(&defValue))
	default:
		log.Fatal().
			Str("question", q.Name).
			Msgf("Unknown prompt type")

		return nil
	}
}
