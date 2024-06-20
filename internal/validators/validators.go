// Package validators contains the validation functions for the application.
// these specificially relate to the prompts/questions asked to the user and
// validating the user input.
package validators

import (
	"errors"
	"fmt"
	"reflect"
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
func NotZero[T comparable](v T) error {
	var zero T
	if v == zero {
		return fmt.Errorf("value is required")
	}

	return nil
}

type lengthType string

var (
	lengthTypeUnknown    lengthType = "unknown"
	lengthTypeCollection lengthType = "collection"
	lengthTypeString     lengthType = "string"
)

// Function to get the length of a value using reflection
func getLength(v interface{}) (lengthType, int) {
	// Get the reflect.Value of the interface
	val := reflect.ValueOf(v)

	// Check if the value is a slice, array, string, map, or channel
	switch val.Kind() {
	case reflect.String:
		return lengthTypeString, val.Len()
	case reflect.Slice, reflect.Array, reflect.Map, reflect.Chan:
		return lengthTypeCollection, val.Len()
	default:
		return lengthTypeUnknown, -1
	}
}

// MinLength validates that the value is atleast m long.
//
// Supported Types:
//   - string
//   - slice
//   - array
//   - map
//   - channel
func MinLength[T any](m int) Validator[T] {
	var zero T
	_, n := getLength(zero)
	if n == -1 {
		panic("unsupported type")
	}

	return func(v T) error {
		t, n := getLength(v)
		if n < m {
			if t == lengthTypeString {
				return fmt.Errorf("input must be greater than %d characters", m)
			} else {
				return fmt.Errorf("must select atleast %d", m)
			}
		}

		return nil
	}
}

// MaxLength validates that the value is at most m long.
//
// Supported Types:
//   - string
//   - slice
//   - array
//   - map
//   - channel
func MaxLength[T any](m int) Validator[T] {
	var zero T
	_, n := getLength(zero)
	if n == -1 {
		panic("unsupported type")
	}

	return func(v T) error {
		t, n := getLength(v)
		if n > m {
			if t == lengthTypeString {
				return fmt.Errorf("input cannot be greater than %d", m)
			} else {
				return fmt.Errorf("cannot select more than %d", m)
			}
		}

		return nil
	}
}

// AtleastOne validates that atleast one selection is required.
func AtleastOne[T any](v []T) error {
	if len(v) == 0 {
		return fmt.Errorf("atleast one selection is required")
	}

	return nil
}

type StringOrStringSlice interface {
	~string | []string
}

func Match[T StringOrStringSlice](re string, msg string) Validator[T] {
	compiled := regexp.MustCompile(re)

	return func(v T) error {
		// Get the reflect.Value of the interface
		val := reflect.ValueOf(v)

		// Check if the value is a slice, array, string, map, or channel
		switch val.Kind() {
		case reflect.String:
			str := val.String()

			if !compiled.MatchString(str) {
				return errors.New(msg)
			}
		case reflect.Slice:
			// cast to string slice
			strSlice := any(v).([]string)
			for _, str := range strSlice {
				if !compiled.MatchString(str) {
					return errors.New(msg)
				}
			}
		default:
			panic("unsupported type")
		}

		return nil
	}
}
