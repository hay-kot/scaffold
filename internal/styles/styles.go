// Package styles contains the shared styles for the terminal UI components.
package styles

import "github.com/charmbracelet/lipgloss"

var (
	Bold  = lipgloss.NewStyle().Bold(true).Render
	Light = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorScaffoldBlueDark)).
		Render
	Base = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorScaffoldBlueSecondary)).
		Render
)

const (
	colorScaffoldBlueLight     = "#5A82E0"
	colorScaffoldBlueDark      = "#758BF9"
	colorScaffoldBlueSecondary = "#059669"
)
