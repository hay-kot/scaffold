// Package validators contains the validation functions for the application.
// these specificially relate to the prompts/questions asked to the user and
// validating the user input.
package validators

import (
	"fmt"
	"reflect"
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

// Function to get the length of a value using reflection
func getLength(v interface{}) int {
	// Get the reflect.Value of the interface
	val := reflect.ValueOf(v)

	// Check if the value is a slice, array, string, map, or channel
	switch val.Kind() {
	case reflect.Slice, reflect.Array, reflect.String, reflect.Map, reflect.Chan:
		return val.Len()
	default:
		return -1 // Return -1 if the value doesn't support length
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
	if getLength(zero) == -1 {
		panic("unsupported type")
	}

	return func(v T) error {
		if getLength(v) < m {
			return fmt.Errorf("value must be at least %d long", m)
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
	if getLength(zero) == -1 {
		panic("unsupported type")
	}

	return func(v T) error {
		if getLength(v) > m {
			return fmt.Errorf("value can be at most %d long", m)
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
