package scaffold

import (
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/internal/styles"
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

func (a *Askable) String() string {
	bldr := strings.Builder{}

	bldr.WriteString("  ")
	bldr.WriteString(styles.Light("?"))
	bldr.WriteString(" ")
	bldr.WriteString(styles.Bold(a.Name))
	bldr.WriteString(" ")

	val := a.Field.GetValue()

	switch v := val.(type) {
	case string:
		if v == "" {
			return ""
		}

		if strings.Contains(v, "\n") {
			bldr.WriteString(styles.Base("|"))

			for _, line := range strings.Split(v, "\n") {
				bldr.WriteString("\n")
				bldr.WriteString("      ")
				bldr.WriteString(styles.Base(line))
			}

			break
		}

		bldr.WriteString(styles.Base(v))
	case []string:
		if len(v) == 0 {
			return ""
		}

		for _, v := range v {
			bldr.WriteString("\n")
			bldr.WriteString("      - ")
			bldr.WriteString(styles.Base(v))
		}
	case bool:
		if v {
			bldr.WriteString(styles.Base("true"))
		} else {
			bldr.WriteString(styles.Base("false"))
		}

	default:
		bldr.WriteString("unknown type, please report this issue to the scaffold maintainer.")
	}

	bldr.WriteString("\n")

	return bldr.String()
}
