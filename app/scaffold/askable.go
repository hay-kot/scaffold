package scaffold

import (
	"github.com/charmbracelet/huh"
	"github.com/hay-kot/scaffold/app/core/engine"
)

type Askable struct {
	Key   string
	Hook  func(vars engine.Vars) error
	Field huh.Field
}

func NewAskable(key string, field huh.Field, fn func(vars engine.Vars) error) *Askable {
	return &Askable{
		Key:   key,
		Field: field,
		Hook:  fn,
	}
}
