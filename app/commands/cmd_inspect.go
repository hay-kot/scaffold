package commands

import (
	"encoding/json"
	"os"

	"github.com/hay-kot/scaffold/app/scaffold"
)

type FlagsInspect struct {
	Path string
}

// InspectOutput is the JSON output format for the inspect command.
type InspectOutput struct {
	Questions []InspectQuestion         `json:"questions"`
	Presets   map[string]map[string]any `json:"presets,omitempty"`
	Computed  map[string]string         `json:"computed,omitempty"`
	Features  []InspectFeature          `json:"features,omitempty"`
	Messages  *InspectMessages          `json:"messages,omitempty"`
}

// InspectQuestion describes a scaffold variable/question.
type InspectQuestion struct {
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Required    bool     `json:"required"`
	Message     string   `json:"message,omitempty"`
	Description string   `json:"description,omitempty"`
	Default     any      `json:"default,omitempty"`
	Options     []string `json:"options,omitempty"`
	Group       string   `json:"group,omitempty"`
}

// InspectFeature describes a scaffold feature toggle.
type InspectFeature struct {
	Value string   `json:"value"`
	Globs []string `json:"globs"`
}

// InspectMessages contains pre/post scaffold messages.
type InspectMessages struct {
	Pre  string `json:"pre,omitempty"`
	Post string `json:"post,omitempty"`
}

func (ctrl *Controller) Inspect(flags FlagsInspect) error {
	ctrl.ready()

	path, err := ctrl.resolve(flags.Path, ".", true, true)
	if err != nil {
		return err
	}

	project, err := scaffold.LoadProject(os.DirFS(path), scaffold.Options{})
	if err != nil {
		return err
	}

	output := InspectOutput{
		Questions: make([]InspectQuestion, len(project.Conf.Questions)),
		Presets:   project.Conf.Presets,
		Computed:  project.Conf.Computed,
	}

	for i, q := range project.Conf.Questions {
		output.Questions[i] = questionToInspect(q)
	}

	if len(project.Conf.Features) > 0 {
		output.Features = make([]InspectFeature, len(project.Conf.Features))
		for i, f := range project.Conf.Features {
			output.Features[i] = InspectFeature{
				Value: f.Value,
				Globs: f.Globs,
			}
		}
	}

	if project.Conf.Messages.Pre != "" || project.Conf.Messages.Post != "" {
		output.Messages = &InspectMessages{
			Pre:  project.Conf.Messages.Pre,
			Post: project.Conf.Messages.Post,
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(output)
}

func questionToInspect(q scaffold.Question) InspectQuestion {
	iq := InspectQuestion{
		Name:     q.Name,
		Required: q.Required || q.Validate.Required,
		Group:    q.Group,
	}

	// Determine type and extract prompt details
	switch {
	case q.Prompt.IsConfirm():
		iq.Type = "bool"
		if q.Prompt.Confirm != nil {
			iq.Message = *q.Prompt.Confirm
		}
	case q.Prompt.IsMultiSelect():
		iq.Type = "[]string"
		if q.Prompt.Message != nil {
			iq.Message = *q.Prompt.Message
		}
		if q.Prompt.Options != nil {
			iq.Options = *q.Prompt.Options
		}
	case q.Prompt.IsSelect():
		iq.Type = "string"
		if q.Prompt.Message != nil {
			iq.Message = *q.Prompt.Message
		}
		if q.Prompt.Options != nil {
			iq.Options = *q.Prompt.Options
		}
	case q.Prompt.IsInputLoop():
		iq.Type = "[]string"
		if q.Prompt.Message != nil {
			iq.Message = *q.Prompt.Message
		}
	default:
		iq.Type = "string"
		if q.Prompt.Message != nil {
			iq.Message = *q.Prompt.Message
		}
	}

	if q.Prompt.Description != nil {
		iq.Description = *q.Prompt.Description
	}

	if q.Prompt.Default != nil {
		iq.Default = q.Prompt.Default
	}

	return iq
}
