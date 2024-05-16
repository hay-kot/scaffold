package styles

import (
	catppuccingo "github.com/catppuccin/go"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type HuhTheme string

var (
	HuhThemeCharm      = HuhTheme("charm")
	HuhThemeDracula    = HuhTheme("dracula")
	HuhThemeBase16     = HuhTheme("base16")
	HuhThemeCatppuccin = HuhTheme("catppuccino")
	HuhThemeScaffold   = HuhTheme("scaffold")
)

// SetGlobalStyles sets the global style reference based on the theme.
func SetGlobalStyles(theme HuhTheme) {
	Theme(theme)
}

// Theme returns a new theme based on the given HuhTheme.
func Theme(theme HuhTheme) *huh.Theme {
	switch theme {
	case HuhThemeCharm:
		// Set Globals
		Light = lipgloss.NewStyle().Foreground(lipgloss.Color("#F780E2")).Render
		Base = lipgloss.NewStyle().Foreground(lipgloss.Color("#7571F9")).Render

		return huh.ThemeCharm()
	case HuhThemeDracula:
		// Set Globals
		Light = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c")).Render
		Base = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4")).Render

		return huh.ThemeDracula()
	case HuhThemeBase16:
		// Set Globals
		Light = lipgloss.NewStyle().Foreground(lipgloss.Color("3")).Render
		Base = lipgloss.NewStyle().Foreground(lipgloss.Color("6")).Render

		return huh.ThemeBase16()
	case HuhThemeCatppuccin:
		// Set Globals
		light := catppuccingo.Latte
		dark := catppuccingo.Mocha

		Light = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: light.Mauve().Hex, Dark: dark.Mauve().Hex}).Render
		Base = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: light.Pink().Hex, Dark: dark.Pink().Hex}).Render

		return huh.ThemeCatppuccin()
	default:
		// Set Globals
		Light = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorScaffoldBlueDark)).
			Render
		Base = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorScaffoldBlueSecondary)).
			Render

		return ThemeScaffold()
	}
}

// ThemeScaffold returns a new theme based on the Charm color scheme.
func ThemeScaffold() *huh.Theme {
	t := huh.ThemeBase()

	var (
		normalFg = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
		indigo   = lipgloss.AdaptiveColor{Light: colorScaffoldBlueLight, Dark: colorScaffoldBlueDark}
		cream    = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
		fuchsia  = lipgloss.Color(colorScaffoldBlueSecondary)
		green    = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
		red      = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Title = t.Focused.Title.Foreground(indigo).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(indigo).Bold(true).MarginBottom(1)
	// t.Focused.Directory = t.Focused.Directory.Foreground(indigo)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"})
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(fuchsia)
	// t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(fuchsia)
	// t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(fuchsia)
	t.Focused.Option = t.Focused.Option.Foreground(normalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(fuchsia)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02CF92", Dark: "#02A877"}).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"}).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(cream).Background(fuchsia)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(normalFg).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(fuchsia)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	// t.Blurred.NextIndicator = lipgloss.NewStyle()
	// t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}
