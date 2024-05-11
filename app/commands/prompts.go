package commands

import (
	"fmt"
	"strings"

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
