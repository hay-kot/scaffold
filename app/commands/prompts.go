package commands

import (
	"fmt"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/charmbracelet/lipgloss"
)

var (
	bold     = lipgloss.NewStyle().Bold(true)
	colorRed = lipgloss.NewStyle().Foreground(lipgloss.Color("#dc2626"))
)

func httpAuthPrompt() (username string, password string, err error) {
	qs := []*survey.Question{
		{
			Name:     "username",
			Prompt:   &survey.Input{Message: "Username:"},
			Validate: survey.Required,
		},
		{
			Name: "password",
			Prompt: &survey.Password{
				Message: "Password/Access Token:",
			},
		},
	}

	answers := struct {
		Username string
		Password string
	}{}

	err = survey.Ask(qs, &answers)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse http auth input: %w", err)
	}

	return answers.Username, answers.Password, nil
}

func didYouMeanPrompt(given, suggestion string) bool {
	bldr := strings.Builder{}

	// Couldn't find a scaffold named:
	//   'foo'
	//
	// Did you mean:
	//   'bar'?
	//
	// [y/n]:

	bldr.WriteString("\n ")
	bldr.WriteString(bold.Render(colorRed.Render("could not find a scaffold named")))
	bldr.WriteString("\n    ")
	bldr.WriteString(given)
	bldr.WriteString("\n\n")
	bldr.WriteString(" ")
	bldr.WriteString(bold.Render("did you mean"))
	bldr.WriteString("\n    ")
	bldr.WriteString(suggestion)
	bldr.WriteString("?\n\n ")
	bldr.WriteString("[y/n]: ")

	out := bldr.String()

	var resp string

	fmt.Print(out)
	fmt.Scanln(&resp)

	return resp == "y"
}
