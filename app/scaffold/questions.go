package scaffold

import (
	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/internal/huhext"
	"github.com/hay-kot/scaffold/internal/validators"
	"github.com/rs/zerolog/log"
)

func QuestionGroupBy(questions []Question) [][]Question {
	grouped := [][]Question{}

outer:
	for _, q := range questions {
		if q.Group == "" {
			grouped = append(grouped, []Question{q})
			continue
		}

		for i, group := range grouped {
			if group[0].Group == q.Group {
				grouped[i] = append(group, q)
				continue outer
			}
		}

		grouped = append(grouped, []Question{q})
	}

	return grouped
}

type Question struct {
	Name     string    `yaml:"name"`
	Group    string    `yaml:"group"`
	Prompt   AnyPrompt `yaml:"prompt"`
	When     string    `yaml:"when"`
	Required bool      `yaml:"required"`
	Validate Validate  `yaml:"validate"`
}

func (q Question) Title() string {
	return unwrap(q.Prompt.Message)
}

func (q Question) Description() string {
	return unwrap(q.Prompt.Desciption)
}

type Validate struct {
	Required  bool `yaml:"required"`
	MinLength int  `yaml:"min"`
	MaxLength int  `yaml:"max"`
	Match     struct {
		Regex   string `yaml:"regex"`
		Message string `yaml:"message"`
	} `yaml:"match"`
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

func (q Question) ToAskable(def any) *Askable {
	switch {
	case q.Prompt.IsMultiSelect():
		defValue := parseDefaultStrings(def, q.Prompt.Default)

		prompt := huh.NewMultiSelect[string]().
			Title(q.Title()).
			Description(q.Description()).
			Options(toHuhOptions(q.Prompt.Options)...).
			Value(&defValue)

		var vals []validators.Validator[[]string]

		if q.Validate.MinLength > 0 {
			vals = append(vals, validators.MinLength[[]string](q.Validate.MinLength))
		} else if q.Required || q.Validate.Required {
			vals = append(vals, validators.AtleastOne[string])
		}

		if q.Validate.MaxLength > 0 {
			vals = append(vals, validators.MaxLength[[]string](q.Validate.MaxLength))
		}

		if len(vals) > 0 {
			prompt.Validate(validators.Combine(vals...))
		}

		return NewAskable(q.Title(), q.Name, prompt, func(vars engine.Vars) error {
			vars[q.Name] = prompt.GetValue().([]string)
			return nil
		})
	case q.Prompt.IsSelect():
		defValue := parseDefaultString(def, q.Prompt.Default)

		prompt := huh.NewSelect[string]().
			Title(q.Title()).
			Description(q.Description()).
			Options(toHuhOptions(q.Prompt.Options)...).
			Value(&defValue)

		if q.Required {
			prompt.Validate(validators.NotZero)
		}

		var vals []validators.Validator[string]

		if q.Validate.MinLength > 0 {
			vals = append(vals, validators.MinLength[string](q.Validate.MinLength))
		} else if q.Required || q.Validate.Required {
			vals = append(vals, validators.NotZero[string])
		}

		if len(vals) > 0 {
			prompt.Validate(validators.Combine(vals...))
		}

		return NewAskable(q.Title(), q.Name, prompt, func(vars engine.Vars) error {
			vars[q.Name] = prompt.GetValue().(string)
			return nil
		})

	case q.Prompt.IsConfirm():
		defValue := parseDefaultBool(def, q.Prompt.Default)
		prompt := huh.NewConfirm().
			Title(q.Title()).
			Description(q.Description()).
			Value(&defValue)

		return NewAskable(q.Title(), q.Name, prompt, func(vars engine.Vars) error {
			vars[q.Name] = prompt.GetValue().(bool)
			return nil
		})
	case q.Prompt.IsInputLoop():
		defValue := parseDefaultStrings(def, q.Prompt.Default)

		prompt := huhext.NewLoopedInput().
			Title(q.Title()).
			Description(q.Description()).
			Value(defValue)

		var vals []validators.Validator[[]string]

		if q.Validate.MinLength > 0 {
			vals = append(vals, validators.MinLength[[]string](q.Validate.MinLength))
		} else if q.Required || q.Validate.Required {
			vals = append(vals, validators.AtleastOne[string])
		}

		if q.Validate.MaxLength > 0 {
			vals = append(vals, validators.MaxLength[[]string](q.Validate.MaxLength))
		}

		if q.Validate.Match.Regex != "" {
			vals = append(vals, validators.Match[[]string](q.Validate.Match.Regex, q.Validate.Match.Message))
		}

		if len(vals) > 0 {
			prompt.Validate(validators.Combine(vals...))
		}

		return NewAskable(q.Title(), q.Name, prompt, func(vars engine.Vars) error {
			vars[q.Name] = prompt.GetValue().([]string)
			return nil
		})

	case q.Prompt.IsTextInput():
		defValue := parseDefaultString(def, q.Prompt.Default)

		prompt := huh.NewText().
			Title(q.Title()).
			Description(q.Description()).
			Value(&defValue)

		var vals []validators.Validator[string]

		if q.Validate.MinLength > 0 {
			vals = append(vals, validators.MinLength[string](q.Validate.MinLength))
		} else if q.Required || q.Validate.Required {
			vals = append(vals, validators.NotZero[string])
		}

		if q.Validate.MaxLength > 0 {
			vals = append(vals, validators.MaxLength[string](q.Validate.MaxLength))
		}

		if q.Validate.Match.Regex != "" {
			vals = append(vals, validators.Match[string](q.Validate.Match.Regex, q.Validate.Match.Message))
		}

		if len(vals) > 0 {
			prompt.Validate(validators.Combine(vals...))
		}

		return NewAskable(q.Title(), q.Name, prompt, func(vars engine.Vars) error {
			vars[q.Name] = prompt.GetValue().(string)
			return nil
		})

	case q.Prompt.IsInput():
		defValue := parseDefaultString(def, q.Prompt.Default)

		prompt := huh.NewInput().
			Title(q.Title()).
			Description(q.Description()).
			Value(&defValue)

		var vals []validators.Validator[string]

		if q.Validate.MinLength > 0 {
			vals = append(vals, validators.MinLength[string](q.Validate.MinLength))
		} else if q.Required || q.Validate.Required {
			vals = append(vals, validators.NotZero[string])
		}

		if q.Validate.MaxLength > 0 {
			vals = append(vals, validators.MaxLength[string](q.Validate.MaxLength))
		}

		if q.Validate.Match.Regex != "" {
			vals = append(vals, validators.Match[string](q.Validate.Match.Regex, q.Validate.Match.Message))
		}

		if len(vals) > 0 {
			prompt.Validate(validators.Combine(vals...))
		}

		return NewAskable(q.Title(), q.Name, prompt, func(vars engine.Vars) error {
			vars[q.Name] = prompt.GetValue().(string)
			return nil
		})

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
