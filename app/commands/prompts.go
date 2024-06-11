package commands

import (
	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/internal/styles"
)

func httpAuthPrompt(theme styles.HuhTheme) (username string, password string, err error) {
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
				EchoMode(huh.EchoModePassword),
		),
	).WithTheme(styles.Theme(theme))

	err = form.Run()
	if err != nil {
		return "", "", err
	}

	return username, password, nil
}

func didYouMeanPrompt(given, suggestion string) bool {
	ok := true

	err := huh.NewConfirm().
		Title("Did You Mean " + suggestion + "?").
		Description("Couldn't find '" + given + "'").
		Value(&ok).
		Run()
	if err != nil {
		return false
	}

	return ok
}
