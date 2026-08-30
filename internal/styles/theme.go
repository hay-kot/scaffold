package styles

import (
	"slices"

	"charm.land/huh/v2"
	"charm.land/lipgloss/v2"
)

type HuhTheme string

var (
	HuhThemeCharm      = HuhTheme("charm")
	HuhThemeDracula    = HuhTheme("dracula")
	HuhThemeBase16     = HuhTheme("base16")
	HuhThemeCatppuccin = HuhTheme("catppuccino")
	HuhThemeScaffold   = HuhTheme("scaffold")
	HuhThemeTokyoNight = HuhTheme("tokyo-night")
)

func (t HuhTheme) String() string {
	return string(t)
}

func (t HuhTheme) IsValid() bool {
	valid := []HuhTheme{
		HuhThemeCharm,
		HuhThemeDracula,
		HuhThemeBase16,
		HuhThemeCatppuccin,
		HuhThemeScaffold,
		HuhThemeTokyoNight,
	}

	return slices.Contains(valid, t)
}

// SetGlobalStyles sets the global style reference based on the theme.
func SetGlobalStyles(theme HuhTheme) {
	Theme(theme)
}

// Theme returns a new theme based on the given HuhTheme.
func Theme(theme HuhTheme) huh.Theme {
	switch theme {
	case HuhThemeCharm:
		Base, Light, Warning = ThemeColorCharm.Compile()

		return huh.ThemeFunc(huh.ThemeCharm)
	case HuhThemeDracula:
		Base, Light, Warning = ThemeColorDracula.Compile()

		return huh.ThemeFunc(huh.ThemeDracula)
	case HuhThemeBase16:
		Base, Light, Warning = ThemeColorsBase16.Compile()

		return huh.ThemeFunc(huh.ThemeBase16)
	case HuhThemeCatppuccin:
		Base, Light, Warning = ThemeColorsCatppuccin.Compile()

		return huh.ThemeFunc(huh.ThemeCatppuccin)
	case HuhThemeTokyoNight:
		Base, Light, Warning = ThemeColorsTokyoNight.Compile()

		return huh.ThemeFunc(ThemeTokyoNight)
	default:
		Base, Light, Warning = ThemeColorsScaffold.Compile()

		return huh.ThemeFunc(ThemeScaffold)
	}
}

// ThemeScaffold returns a new theme based on the Charm color scheme.
func ThemeScaffold(isDark bool) *huh.Styles {
	t := huh.ThemeBase(isDark)
	lightDark := lipgloss.LightDark(isDark)

	var (
		normalFg  = lightDark(lipgloss.Color("235"), lipgloss.Color("252"))
		primary   = lightDark(lipgloss.Color(ThemeColorsScaffold.Base), lipgloss.Color(ThemeColorsScaffold.BaseDark))
		cream     = lipgloss.Color("#FFFDF5")
		secondary = lipgloss.Color(ThemeColorsScaffold.Light)
		green     = lightDark(lipgloss.Color("#02BA84"), lipgloss.Color("#02BF87"))
		red       = lightDark(lipgloss.Color("#FF4672"), lipgloss.Color("#ED567A"))
		subtle    = lightDark(lipgloss.Color(""), lipgloss.Color("243"))
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Title = t.Focused.Title.Foreground(primary).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(primary).Bold(true).MarginBottom(1)
	// t.Focused.Directory = t.Focused.Directory.Foreground(secondary)
	t.Focused.Description = t.Focused.Description.Foreground(subtle)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(secondary)
	// t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(primary)
	// t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(primary)
	t.Focused.Option = t.Focused.Option.Foreground(normalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(secondary)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().
		Foreground(lightDark(lipgloss.Color("#02CF92"), lipgloss.Color("#02A877"))).
		SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(subtle).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(cream).Background(secondary)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.
		Foreground(normalFg).
		Background(lightDark(lipgloss.Color("252"), lipgloss.Color("237")))

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.
		Foreground(lightDark(lipgloss.Color("248"), lipgloss.Color("238")))
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(secondary)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	// t.Blurred.NextIndicator = lipgloss.NewStyle()
	// t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}

// ThemeTokyoNight returns a new theme based on the Tokyo Night color scheme.
func ThemeTokyoNight(isDark bool) *huh.Styles {
	t := huh.ThemeBase(isDark)
	lightDark := lipgloss.LightDark(isDark)

	var (
		normalFg   = lightDark(lipgloss.Color("235"), lipgloss.Color("#c0caf5"))
		primary    = lightDark(lipgloss.Color(ThemeColorsTokyoNight.Base), lipgloss.Color(ThemeColorsTokyoNight.BaseDark))
		background = lipgloss.Color("#1a1b26")
		secondary  = lipgloss.Color(ThemeColorsTokyoNight.Light)
		blue       = lipgloss.Color("#7aa2f7")
		green      = lipgloss.Color("#9ece6a")
		red        = lipgloss.Color("#f7768e")
		cyan       = lipgloss.Color("#7dcfff")
		subtle     = lightDark(lipgloss.Color(""), lipgloss.Color("#565f89"))
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("#414868"))
	t.Focused.Title = t.Focused.Title.Foreground(primary).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(primary).Bold(true).MarginBottom(1)
	t.Focused.Description = t.Focused.Description.Foreground(subtle)
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(secondary)
	t.Focused.Option = t.Focused.Option.Foreground(normalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(secondary)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(cyan)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(green).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(subtle).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(background).Background(blue)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.
		Foreground(normalFg).
		Background(lightDark(lipgloss.Color("252"), lipgloss.Color("#292e42")))

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(cyan)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.
		Foreground(lightDark(lipgloss.Color("248"), lipgloss.Color("#414868")))
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(secondary)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())

	return t
}
