// Package rule provides a straightforward and flexible way to handle rule-based
// logic.
package rule

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type Rule int

const (
	Unset Rule = iota
	Yes
	No
	Prompt
)

func NewFromString(s string) (Rule, error) {
	switch s {
	case "yes":
		return Yes, nil
	case "no":
		return No, nil
	case "prompt":
		return Prompt, nil
	}

	return Unset, fmt.Errorf("invalid rule: %v", s)
}

func (r Rule) String() string {
	switch r {
	case Unset:
		return "unset"
	case Yes:
		return "yes"
	case No:
		return "no"
	case Prompt:
		return "prompt"
	}
	panic("invalid rule")
}

func (r *Rule) UnmarshalYAML(node *yaml.Node) error {
	var asBool bool
	var asString string

	switch {
	case node.Decode(&asBool) == nil:
		if asBool {
			*r = Yes
		} else {
			*r = No
		}
		return nil
	case node.Decode(&asString) == nil:
		rule, err := NewFromString(asString)
		if err != nil {
			return err
		}
		*r = rule
		return nil
	default:
		return fmt.Errorf("invalid rule: %v", node.Value)
	}
}
