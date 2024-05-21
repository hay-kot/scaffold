// Package styles contains the shared styles for the terminal UI components.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	Bold        = lipgloss.NewStyle().Bold(true).Render
	Underline   = lipgloss.NewStyle().Underline(true).Render
	Base, Light = ThemeColorsScaffold.Compile()
)
