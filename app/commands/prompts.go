package commands

import (
	"errors"
	"fmt"

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

func scaffoldPickerPrompt(localScaffolds []string, systemScaffolds []string, theme styles.HuhTheme) (string, error) {
	if len(localScaffolds) == 0 && len(systemScaffolds) == 0 {
		return "", errors.New("no scaffolds available")
	}

	options := make([]huh.Option[string], 0, len(localScaffolds)+len(systemScaffolds))

	for _, name := range localScaffolds {
		options = append(options, huh.NewOption(fmt.Sprintf("[local]  %s", name), name))
	}

	for _, name := range systemScaffolds {
		options = append(options, huh.NewOption(fmt.Sprintf("[system] %s", name), "https://"+name))
	}

	var selected string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a scaffold").
				Options(options...).
				Filtering(true).
				Value(&selected),
		),
	).WithTheme(styles.Theme(theme))

	if err := form.Run(); err != nil {
		return "", err
	}

	return selected, nil
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
