package scaffold

import (
	"cmp"
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/rs/zerolog/log"
)

type Question struct {
	Name     string    `yaml:"name"`
	Prompt   AnyPrompt `yaml:"prompt"`
	When     string    `yaml:"when"`
	Required bool      `yaml:"required"`
}

func (q Question) Title() string {
	return unwrap(q.Prompt.Message)
}

func (q Question) Description() string {
	return unwrap(q.Prompt.Desciption)
}

type AnyPrompt struct {
	Message    *string   `yaml:"message"`
	Desciption *string   `yaml:"description"`
	Loop       bool      `yaml:"loop"`
	Default    any       `yaml:"default"`
	Confirm    *string   `yaml:"confirm"`
	Multi      bool      `yaml:"multi"`
	Options    *[]string `yaml:"options"`
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

func (p AnyPrompt) IsTextInput() bool {
	return p.IsInput() && !p.Loop && p.Multi
}

func (p AnyPrompt) IsMultiSelect() bool {
	return p.IsSelect() && p.Multi
}

func (q Question) ToAskable(defaults engine.Vars) Askable {
	def := defaults[q.Name]

	switch {
	case q.Prompt.IsMultiSelect():
		defValue := parseDefaultStrings(def, q.Prompt.Default)

		prompt := huh.NewMultiSelect[string]().
			Title(q.Title()).
			Description(q.Description()).
			Options(toHuhOptions(q.Prompt.Options)...).
			Value(&defValue)

		if q.Required {
			prompt.Validate(validateAtleastOne)
		}

		return HuhToAskable[[]string](q.Name, prompt)

	case q.Prompt.IsSelect():
		defValue := parseDefaultString(def, q.Prompt.Default)

		prompt := huh.NewSelect[string]().
			Title(q.Title()).
			Description(q.Description()).
			Options(toHuhOptions(q.Prompt.Options)...).
			Value(&defValue)

		if q.Required {
			prompt.Validate(validateNotZero)
		}

		return HuhToAskable[string](q.Name, prompt)
	case q.Prompt.IsConfirm():
		defValue := parseDefaultBool(def, q.Prompt.Default)
		return HuhToAskable[bool](q.Name, huh.NewConfirm().
			Title(q.Title()).
			Description(q.Description()).
			Value(&defValue))
	case q.Prompt.IsInputLoop():
		var out []string
		return AskableFunc(func(vars engine.Vars) error {
			for {
				ref := ""

				desc := cmp.Or(q.Description(), "Press enter to finish")

				ask := huh.NewInput().
					Title(q.Title()).
					Description(desc).
					Value(&ref)

				if q.Required {
					ask.Validate(func(s string) error {
						if s == "" && len(out) == 0 {
							return fmt.Errorf("atleast one value is required")
						}

						return nil
					})
				}

				err := ask.Run()
				if err != nil {
					return nil
				}

				fmt.Println(ask.View())
				fmt.Print("\n")

				if ref == "" {
					break
				}

				out = append(out, ref)
			}

			vars[q.Name] = out
			return nil
		})

	case q.Prompt.IsTextInput():
		defValue := parseDefaultString(def, q.Prompt.Default)
		desc := cmp.Or(q.Description(), "ctrl+j for new line")

		prompt := huh.NewText().
			Title(q.Title()).
			Description(desc).
			Value(&defValue)

		if q.Required {
			prompt.Validate(validateNotZero)
		}

		return HuhToAskable[string](q.Name, prompt)
	case q.Prompt.IsInput():
		defValue := parseDefaultString(def, q.Prompt.Default)

		prompt := huh.NewInput().
			Title(q.Title()).
			Description(q.Description()).
			Value(&defValue)

		if q.Required {
			prompt.Validate(validateNotZero)
		}

		return HuhToAskable[string](q.Name, prompt)
	default:
		log.Fatal().
			Str("question", q.Name).
			Msgf("Unknown prompt type")

		return nil
	}
}

func toHuhOptions(opts *[]string) []huh.Option[string] {
	out := make([]huh.Option[string], len(*opts))
	for i, opt := range *opts {
		out[i] = huh.NewOption(opt, opt)
	}
	return out
}

func unwrap[T any](v *T) T {
	if v == nil {
		var zero T
		return zero
	}

	return *v
}

func validateNotZero[T comparable](v T) error {
	var zero T
	if v == zero {
		return fmt.Errorf("value is required")
	}

	return nil
}

func validateAtleastOne[T any](v []T) error {
	if len(v) == 0 {
		return fmt.Errorf("atleast one selection is required")
	}

	return nil
}
