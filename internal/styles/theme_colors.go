package styles

import (
	catppuccingo "github.com/catppuccin/go"
	"github.com/charmbracelet/lipgloss"
)

var (
	ThemeColorCharm     = &ThemeColors{Base: "#7571F9", Light: "#F780E2"}
	ThemeColorDracula   = &ThemeColors{Base: "#6272a4", Light: "#F1FA8C"}
	ThemeColorsBase16   = &ThemeColors{Base: "6", Light: "3"}
	ThemeColorsScaffold = &ThemeColors{
		Base:     "#5A82E0",
		BaseDark: "#758BF9",
		Light:    "#059669",
	}
	ThemeColorsCatppuccin = &ThemeColors{
		Base:      catppuccingo.Latte.Mauve().Hex,
		BaseDark:  catppuccingo.Mocha.Mauve().Hex,
		Light:     catppuccingo.Latte.Pink().Hex,
		LightDark: catppuccingo.Mocha.Pink().Hex,
	}
)

type ThemeColors struct {
	Base      string
	BaseDark  string
	Light     string
	LightDark string

	compiled bool
	base     lipgloss.Style
	light    lipgloss.Style
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
}

func (t *ThemeColors) Compile() (base RenderFunc, light RenderFunc) {
	t.compile()
	return t.base.Render, t.light.Render
}

func (t *ThemeColors) BaseFn(string ...string) string {
	t.compile()
	return t.base.Render(string...)
}

func (t *ThemeColors) LightFn(string ...string) string {
	t.compile()
	return t.light.Render(string...)
}
