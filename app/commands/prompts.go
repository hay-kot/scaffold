package commands

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

var (
	bold     = lipgloss.NewStyle().Bold(true)
	colorRed = lipgloss.NewStyle().Foreground(lipgloss.Color("#dc2626"))
)

func httpAuthPrompt() (username string, password string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Description("Enter your username").
				Value(&username),
			huh.NewInput().
				Title("Password").
				Description("Enter your password (or token)").
				Value(&password).
				Password(true),
		),
	)

	err = form.Run()
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func didYouMeanPrompt(given, suggestion string) bool {
	// Couldn't find a scaffold named:
	//   'foo'
	//
	// Did you mean:
	//   'bar'?
	//
	// [y/n]:

	ok := true

	fmt.Print("\n")
	err := huh.NewConfirm().Title("Couldn't Find '" + given + "'").
		Description("Did you mean: " + suggestion + "?").
		Value(&ok).
		Run()
	if err != nil {
		return false
	}

	return ok
}
