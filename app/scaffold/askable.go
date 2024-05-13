package scaffold

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/scaffold/app/core/engine"
)

type Askable struct {
	Name  string
	Key   string
	Hook  func(vars engine.Vars) error
	Field huh.Field
}

func NewAskable(name string, key string, field huh.Field, fn func(vars engine.Vars) error) *Askable {
	return &Askable{
		Name:  name,
		Key:   key,
		Field: field,
		Hook:  fn,
	}
}

var (
	bold  = lipgloss.NewStyle().Bold(true).Render
	light = lipgloss.NewStyle().Foreground(lipgloss.Color("#FF83D5")).Render
	base  = lipgloss.NewStyle().Foreground(lipgloss.Color("#6C74FD")).Render
)

func (a *Askable) String() string {
	bldr := strings.Builder{}

	bldr.WriteString("  ")
	bldr.WriteString(light("?"))
	bldr.WriteString(" ")
	bldr.WriteString(bold(a.Name))
	bldr.WriteString(" ")

	val := a.Field.GetValue()

	switch v := val.(type) {
	case string:
		if v == "" {
			return ""
		}

		if strings.Contains(v, "\n") {
			bldr.WriteString(base("|"))

			for _, line := range strings.Split(v, "\n") {
				bldr.WriteString("\n")
				bldr.WriteString("      ")
				bldr.WriteString(base(line))
			}

			break
		}

		bldr.WriteString(base(v))
	case []string:
		if len(v) == 0 {
			return ""
		}

		for _, v := range v {
			bldr.WriteString("\n")
			bldr.WriteString("      - ")
			bldr.WriteString(base(v))
		}
	case bool:
		if v {
			bldr.WriteString(base("true"))
		} else {
			bldr.WriteString(base("false"))
		}

	default:
		bldr.WriteString("unknown type, please report this issue to the scaffold maintainer.")
	}

	bldr.WriteString("\n")

	return bldr.String()
}
