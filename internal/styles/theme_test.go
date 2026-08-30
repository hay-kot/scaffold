package styles_test

import (
	"testing"

	"charm.land/huh/v2"
	"github.com/hay-kot/scaffold/internal/styles"
)

// TestThemeBuildsForm checks that every theme produces styles and that a form
// renders with them. Huh v2 resolves a theme through an interface and per
// light/dark call, so a broken theme fails at run time, not at compile time.
func TestThemeBuildsForm(t *testing.T) {
	names := []styles.HuhTheme{
		styles.HuhThemeCharm,
		styles.HuhThemeDracula,
		styles.HuhThemeBase16,
		styles.HuhThemeCatppuccin,
		styles.HuhThemeScaffold,
		styles.HuhThemeTokyoNight,
		styles.HuhTheme("unknown"),
	}

	for _, name := range names {
		t.Run(name.String(), func(t *testing.T) {
			theme := styles.Theme(name)

			for _, isDark := range []bool{true, false} {
				if theme.Theme(isDark) == nil {
					t.Fatalf("nil styles for isDark=%v", isDark)
				}
			}

			form := huh.NewForm(huh.NewGroup(
				huh.NewInput().Title("Name"),
				huh.NewSelect[string]().Title("Pick").Options(huh.NewOptions("a", "b")...),
				huh.NewConfirm().Title("Sure?"),
			)).WithTheme(theme)

			if form.View() == "" {
				t.Fatal("form rendered an empty view")
			}
		})
	}
}
