package commands

import (
	"errors"
	"fmt"
	"sort"

	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/internal/styles"
)

func httpAuthPrompt(pkgurl string, theme styles.HuhTheme) (username string, password string, err error) {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Username").
				Description(fmt.Sprintf("Authentication required for %s", pkgurl)).
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

func scaffoldPickerPrompt(aliases map[string]string, localScaffolds []string, systemScaffolds []string, theme styles.HuhTheme) (string, error) {
	if len(aliases) == 0 && len(localScaffolds) == 0 && len(systemScaffolds) == 0 {
		return "", errors.New("no scaffolds available, run 'scaffold update' to fetch scaffolds or 'scaffold init' to create local scaffolds")
	}

	options := make([]huh.Option[string], 0, len(aliases)+len(localScaffolds)+len(systemScaffolds))

	if len(aliases) > 0 {
		names := make([]string, 0, len(aliases))
		for name := range aliases {
			names = append(names, name)
		}
		sort.Strings(names)

		for _, name := range names {
			options = append(options, huh.NewOption(fmt.Sprintf("[alias]  %s → %s", name, aliases[name]), name))
		}
	}

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

func didYouMeanPrompt(given, suggestion string, isSystem bool) bool {
	ok := true

	source := "local"
	if isSystem {
		source = "system"
	}

	err := huh.NewConfirm().
		Title(fmt.Sprintf("Did You Mean %s (%s)?", suggestion, source)).
		Description("Couldn't find '" + given + "'").
		Value(&ok).
		Run()
	if err != nil {
		return false
	}

	return ok
}
