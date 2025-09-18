// Package styles contains the shared styles for the terminal UI components.
package styles

import (
	"github.com/charmbracelet/lipgloss"
)

const (
	Check  = "✔"
	Cross  = "✘"
	Git    = "\uf02a2" // 󰊢
	Folder = ""
	Dot    = "•"
)

const (
	ColorSuccess = "#22c55e"
	ColorError   = "#d75f6b"
)

var (
	Bold        = lipgloss.NewStyle().Bold(true).Render
	Padding     = lipgloss.NewStyle().PaddingLeft(1).Render
	Underline   = lipgloss.NewStyle().Underline(true).Render
	Base, Light = ThemeColorsScaffold.Compile()

	Error   = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorError)).Render
	Success = lipgloss.NewStyle().Foreground(lipgloss.Color(ColorSuccess)).PaddingLeft(1).Render
	Subtle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#a3a3a3")).PaddingLeft(1).Render
)
