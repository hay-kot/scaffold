package styles

import (
	catppuccingo "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var (
	ThemeColorCharm = &ThemeColors{
		Base:    "#7571F9",
		Light:   "#F780E2",
		Warning: "#FFBD5A",
	}
	ThemeColorDracula = &ThemeColors{
		Base:    "#6272a4",
		Light:   "#F1FA8C",
		Warning: "#FFB86C",
	}
	ThemeColorsBase16 = &ThemeColors{
		Base:    "6",
		Light:   "3",
		Warning: "11",
	}
	ThemeColorsScaffold = &ThemeColors{
		Base:     "#5A82E0",
		BaseDark: "#758BF9",
		Light:    "#059669",
		Warning:  "#F59E0B",
	}
	ThemeColorsCatppuccin = &ThemeColors{
		Base:        catppuccingo.Latte.Mauve().Hex,
		BaseDark:    catppuccingo.Mocha.Mauve().Hex,
		Light:       catppuccingo.Latte.Pink().Hex,
		LightDark:   catppuccingo.Mocha.Pink().Hex,
		Warning:     catppuccingo.Latte.Peach().Hex,
		WarningDark: catppuccingo.Mocha.Peach().Hex,
	}
	ThemeColorsTokyoNight = &ThemeColors{
		Base:        "#7aa2f7",
		BaseDark:    "#7aa2f7",
		Light:       "#9ece6a",
		Warning:     "#ff9e64",
		WarningDark: "#ff9e64",
	}
)

type ThemeColors struct {
	Base        string
	BaseDark    string
	Light       string
	LightDark   string
	Warning     string
	WarningDark string

	compiled bool
	base     lipgloss.Style
	light    lipgloss.Style
	warning  lipgloss.Style
}

type RenderFunc func(string ...string) string

func (t *ThemeColors) compile() {
	if t.compiled {
		return
	}

	if t.BaseDark == "" {
		t.BaseDark = t.Base
	}
	t.base = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: t.Base,
		Dark:  t.BaseDark,
	})

	if t.LightDark == "" {
		t.LightDark = t.Light
	}
	t.light = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: t.Light,
		Dark:  t.LightDark,
	})

	if t.WarningDark == "" {
		t.WarningDark = t.Warning
	}
	t.warning = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Light: t.Warning,
		Dark:  t.WarningDark,
	})
}

func (t *ThemeColors) Compile() (base RenderFunc, light RenderFunc, warning RenderFunc) {
	t.compile()
	return t.base.Render, t.light.Render, t.warning.Render
}

func (t *ThemeColors) Styles() (base lipgloss.Style, light lipgloss.Style, warning lipgloss.Style) {
	t.compile()
	return t.base, t.light, t.warning
}

func (t *ThemeColors) BaseFn(string ...string) string {
	t.compile()
	return t.base.Render(string...)
}

func (t *ThemeColors) LightFn(string ...string) string {
	t.compile()
	return t.light.Render(string...)
}

func (t *ThemeColors) WarningFn(string ...string) string {
	t.compile()
	return t.warning.Render(string...)
}
