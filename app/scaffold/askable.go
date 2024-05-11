package scaffold

import (
	"fmt"

	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/engine"
)

// Askable is an interface for types that can ask questions.
// Askable types are used to prompt the user for input.
type Askable interface {
	Ask(vars engine.Vars) error
}

type AskableFunc func(vars engine.Vars) error

func (a AskableFunc) Ask(vars engine.Vars) error {
	return a(vars)
}

// HuhToAskable converts a huh.Field into an Askable type. The name is the key
// to store the value in the vars map. Note that T must be the output value
// to the huh.Field. If the T is not the same as the huh.Field output, the
// function will panic as we assert the type.
func HuhToAskable[T any](name string, h huh.Field) Askable {
	return AskableFunc(func(vars engine.Vars) error {
		err := h.Run()
		if err != nil {
			return err
		}

		fmt.Println(h.View())
		fmt.Print("\n")

		vars[name] = h.GetValue().(T)
		return nil
	})
}
