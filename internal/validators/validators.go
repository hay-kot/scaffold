// Package validators contains the validation functions for the application.
// these specificially relate to the prompts/questions asked to the user and
// validating the user input.
package validators

import (
	"errors"
	"fmt"
	"regexp"
)

type Validator[T any] func(v T) error

func Combine[T any](validators ...Validator[T]) Validator[T] {
	return func(v T) error {
		for _, validate := range validators {
			if err := validate(v); err != nil {
				return err
			}
		}

		return nil
	}
}

// NotZero validates that the value is not zero.
// It returns an error if the value is zero.
func NotZero[T Validatable](v T) error {
	handler := &validatehandler{
		strfn: func(str string) error {
			if str == "" {
				return errors.New("input cannot be empty")
			}

			return nil
		},
		slicefn: func(slice []string) error {
			if len(slice) == 0 {
				return errors.New("input cannot be empty")
			}

			return nil
		},
	}

	return handler.validate(v)
}

// MinLength validates that the value is at least m long.
func MinLength[T Validatable](m int) Validator[T] {
	handler := &validatehandler{
		strfn: func(str string) error {
			if len(str) < m {
				return fmt.Errorf("input must be at least %d characters", m)
			}

			return nil
		},
		slicefn: func(slice []string) error {
			if len(slice) < m {
				return fmt.Errorf("must select at least %d", m)
			}

			return nil
		},
	}

	return func(v T) error { return handler.validate(v) }
}

// MaxLength validates that the value is at most m long.
func MaxLength[T any](m int) Validator[T] {
	handler := &validatehandler{
		strfn: func(str string) error {
			if len(str) > m {
				return fmt.Errorf("input must be at most %d characters", m)
			}

			return nil
		},
		slicefn: func(slice []string) error {
			if len(slice) > m {
				return fmt.Errorf("must select at most %d", m)
			}

			return nil
		},
	}

	return func(v T) error { return handler.validate(v) }
}

func Match[T Validatable](re string, msg string) Validator[T] {
	compiled := regexp.MustCompile(re)

	handler := &validatehandler{
		strfn: func(str string) error {
			if !compiled.MatchString(str) {
				return errors.New(msg)
			}

			return nil
		},
		slicefn: func(slice []string) error {
			for _, str := range slice {
				if !compiled.MatchString(str) {
					return errors.New(msg)
				}
			}

			return nil
		},
	}

	return func(v T) error { return handler.validate(v) }
}
