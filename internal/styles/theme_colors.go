package styles

import (
	"os"
	"sync"

	"charm.land/lipgloss/v2"
	catppuccingo "github.com/catppuccin/go"
)

// hasDarkBackground reports whether the terminal has a dark background.
//
// Lip Gloss v2 removed AdaptiveColor, which picked its light or dark variant at
// render time. A style must now be built with the answer already known. The
// query puts the terminal in raw mode and waits for an OSC 11 answer, so ask at
// most once, and only when a style is first built.
var hasDarkBackground = sync.OnceValue(func() bool {
	return lipgloss.HasDarkBackground(os.Stdin, os.Stdout)
})

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

	once    sync.Once
	base    lipgloss.Style
	light   lipgloss.Style
	warning lipgloss.Style
}

type RenderFunc func(string ...string) string

func (t *ThemeColors) compile() {
	t.once.Do(func() {
		lightDark := lipgloss.LightDark(hasDarkBackground())

		if t.BaseDark == "" {
			t.BaseDark = t.Base
		}
		t.base = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color(t.Base), lipgloss.Color(t.BaseDark)))

		if t.LightDark == "" {
			t.LightDark = t.Light
		}
		t.light = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color(t.Light), lipgloss.Color(t.LightDark)))

		if t.WarningDark == "" {
			t.WarningDark = t.Warning
		}
		t.warning = lipgloss.NewStyle().
			Foreground(lightDark(lipgloss.Color(t.Warning), lipgloss.Color(t.WarningDark)))
	})
}

// Compile returns render functions for the theme colors. It hands back the
// method values rather than the compiled styles so that the background query in
// compile is deferred to the first render. Commands that print nothing styled,
// such as --help, then never touch the terminal state.
func (t *ThemeColors) Compile() (base RenderFunc, light RenderFunc, warning RenderFunc) {
	return t.BaseFn, t.LightFn, t.WarningFn
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
