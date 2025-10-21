package styles

import (
	"slices"

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
	}

	return slices.Contains(valid, t)
}

// SetGlobalStyles sets the global style reference based on the theme.
func SetGlobalStyles(theme HuhTheme) {
	Theme(theme)
}

// Theme returns a new theme based on the given HuhTheme.
func Theme(theme HuhTheme) *huh.Theme {
	switch theme {
	case HuhThemeCharm:
		Base, Light = ThemeColorCharm.Compile()

		return huh.ThemeCharm()
	case HuhThemeDracula:
		Base, Light = ThemeColorDracula.Compile()

		return huh.ThemeDracula()
	case HuhThemeBase16:
		Base, Light = ThemeColorsBase16.Compile()

		return huh.ThemeBase16()
	case HuhThemeCatppuccin:
		Base, Light = ThemeColorsCatppuccin.Compile()

		return huh.ThemeCatppuccin()
	default:
		Base, Light = ThemeColorsScaffold.Compile()

		return ThemeScaffold()
	}
}

// ThemeScaffold returns a new theme based on the Charm color scheme.
func ThemeScaffold() *huh.Theme {
	t := huh.ThemeBase()

	var (
		normalFg  = lipgloss.AdaptiveColor{Light: "235", Dark: "252"}
		primary   = lipgloss.AdaptiveColor{Light: ThemeColorsScaffold.Base, Dark: ThemeColorsScaffold.BaseDark}
		cream     = lipgloss.AdaptiveColor{Light: "#FFFDF5", Dark: "#FFFDF5"}
		secondary = lipgloss.Color(ThemeColorsScaffold.Light)
		green     = lipgloss.AdaptiveColor{Light: "#02BA84", Dark: "#02BF87"}
		red       = lipgloss.AdaptiveColor{Light: "#FF4672", Dark: "#ED567A"}
	)

	t.Focused.Base = t.Focused.Base.BorderForeground(lipgloss.Color("238"))
	t.Focused.Title = t.Focused.Title.Foreground(primary).Bold(true)
	t.Focused.NoteTitle = t.Focused.NoteTitle.Foreground(primary).Bold(true).MarginBottom(1)
	// t.Focused.Directory = t.Focused.Directory.Foreground(secondary)
	t.Focused.Description = t.Focused.Description.Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"})
	t.Focused.ErrorIndicator = t.Focused.ErrorIndicator.Foreground(red)
	t.Focused.ErrorMessage = t.Focused.ErrorMessage.Foreground(red)
	t.Focused.SelectSelector = t.Focused.SelectSelector.Foreground(secondary)
	// t.Focused.NextIndicator = t.Focused.NextIndicator.Foreground(primary)
	// t.Focused.PrevIndicator = t.Focused.PrevIndicator.Foreground(primary)
	t.Focused.Option = t.Focused.Option.Foreground(normalFg)
	t.Focused.MultiSelectSelector = t.Focused.MultiSelectSelector.Foreground(secondary)
	t.Focused.SelectedOption = t.Focused.SelectedOption.Foreground(green)
	t.Focused.SelectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "#02CF92", Dark: "#02A877"}).SetString("✓ ")
	t.Focused.UnselectedPrefix = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "", Dark: "243"}).SetString("• ")
	t.Focused.UnselectedOption = t.Focused.UnselectedOption.Foreground(normalFg)
	t.Focused.FocusedButton = t.Focused.FocusedButton.Foreground(cream).Background(secondary)
	t.Focused.Next = t.Focused.FocusedButton
	t.Focused.BlurredButton = t.Focused.BlurredButton.Foreground(normalFg).Background(lipgloss.AdaptiveColor{Light: "252", Dark: "237"})

	t.Focused.TextInput.Cursor = t.Focused.TextInput.Cursor.Foreground(green)
	t.Focused.TextInput.Placeholder = t.Focused.TextInput.Placeholder.Foreground(lipgloss.AdaptiveColor{Light: "248", Dark: "238"})
	t.Focused.TextInput.Prompt = t.Focused.TextInput.Prompt.Foreground(secondary)

	t.Blurred = t.Focused
	t.Blurred.Base = t.Focused.Base.BorderStyle(lipgloss.HiddenBorder())
	// t.Blurred.NextIndicator = lipgloss.NewStyle()
	// t.Blurred.PrevIndicator = lipgloss.NewStyle()

	return t
}
