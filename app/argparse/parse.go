package argparse

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

// Parse parses command-line arguments in the format key[:type]=value
// and returns a map of properly typed values.
//
// Supported formats:
//   - key=value                    (string, default)
//   - key:string=value             (explicit string)
//   - key:str=value                (string shorthand)
//   - key:int=value                (integer)
//   - key:float=value              (float64)
//   - key:float32=value            (float32)
//   - key:float64=value            (float64)
//   - key:bool=value               (boolean)
//   - key:[]string=val1,val2       (string slice)
//   - key:[]str=val1,val2          (string slice shorthand)
//   - key:[]int=1,2,3              (int slice)
//   - key:[]float=1.1,2.2          (float64 slice)
//   - key:[]float32=1.1,2.2        (float32 slice)
//   - key:[]float64=1.1,2.2        (float64 slice)
//   - key:[]bool=true,false        (bool slice)
//   - key:json={"foo":"bar"}       (JSON value)
//
// For slices, commas can be escaped with backslash: key:[]string=has\,comma,normal
func Parse(args []string) (map[string]any, error) {
	vars := make(map[string]any, len(args))

	for _, arg := range args {
		key, typeHint, value, err := parseArgument(arg)
		if err != nil {
			return nil, fmt.Errorf("invalid argument %q: %w", arg, err)
		}

		parsedValue, err := parseValue(value, typeHint)
		if err != nil {
			return nil, fmt.Errorf("failed to parse %q as %s: %w", key, typeHint, err)
		}

		vars[key] = parsedValue
	}

	return vars, nil
}

// parseArgument splits an argument into key, type hint, and value components
func parseArgument(arg string) (key, typeHint, value string, err error) {
	// Find the first = sign to split key[:type] from value
	eqIdx := strings.Index(arg, "=")
	if eqIdx == -1 {
		return "", "", "", fmt.Errorf("missing '=' in argument")
	}

	keyPart := arg[:eqIdx]
	value = arg[eqIdx+1:] // Everything after first = is the value

	if keyPart == "" {
		return "", "", "", fmt.Errorf("empty key")
	}

	// Check if there's a type hint
	if colonIdx := strings.Index(keyPart, ":"); colonIdx != -1 {
		key = keyPart[:colonIdx]
		typeHint = keyPart[colonIdx+1:]

		if key == "" {
			return "", "", "", fmt.Errorf("empty key before type hint")
		}
		if typeHint == "" {
			return "", "", "", fmt.Errorf("empty type hint after ':'")
		}
	} else {
		key = keyPart
		typeHint = "string" // default type
	}

	return key, typeHint, value, nil
}

// parseValue converts a string value to the appropriate type based on the type hint
func parseValue(value string, typeHint string) (any, error) {
	switch typeHint {
	case "string", "str":
		return value, nil

	case "int":
		return strconv.Atoi(value)

	case "float", "float64":
		return strconv.ParseFloat(value, 64)

	case "float32":
		f, err := strconv.ParseFloat(value, 32)
		if err != nil {
			return nil, err
		}
		return float32(f), nil

	case "bool":
		return strconv.ParseBool(value)

	case "[]string", "[]str":
		if value == "" {
			return []string{}, nil
		}
		parts := splitEscaped(value, ',')
		return parts, nil

	case "[]int":
		if value == "" {
			return []int{}, nil
		}
		parts := splitEscaped(value, ',')
		result := make([]int, len(parts))
		for i, part := range parts {
			val, err := strconv.Atoi(part)
			if err != nil {
				return nil, fmt.Errorf("invalid int %q at position %d", part, i)
			}
			result[i] = val
		}
		return result, nil

	case "[]float", "[]float64":
		if value == "" {
			return []float64{}, nil
		}
		parts := splitEscaped(value, ',')
		result := make([]float64, len(parts))
		for i, part := range parts {
			val, err := strconv.ParseFloat(part, 64)
			if err != nil {
				return nil, fmt.Errorf("invalid float %q at position %d", part, i)
			}
			result[i] = val
		}
		return result, nil

	case "[]float32":
		if value == "" {
			return []float32{}, nil
		}
		parts := splitEscaped(value, ',')
		result := make([]float32, len(parts))
		for i, part := range parts {
			val, err := strconv.ParseFloat(part, 32)
			if err != nil {
				return nil, fmt.Errorf("invalid float32 %q at position %d", part, i)
			}
			result[i] = float32(val)
		}
		return result, nil

	case "[]bool":
		if value == "" {
			return []bool{}, nil
		}
		parts := splitEscaped(value, ',')
		result := make([]bool, len(parts))
		for i, part := range parts {
			val, err := strconv.ParseBool(part)
			if err != nil {
				return nil, fmt.Errorf("invalid bool %q at position %d", part, i)
			}
			result[i] = val
		}
		return result, nil

	case "json":
		var result any
		if err := json.Unmarshal([]byte(value), &result); err != nil {
			return nil, fmt.Errorf("invalid JSON: %w", err)
		}
		return result, nil

	default:
		return nil, fmt.Errorf("unknown type %q", typeHint)
	}
}

// splitEscaped splits a string by the given separator, respecting backslash escapes
func splitEscaped(s string, sep rune) []string {
	if s == "" {
		return []string{""}
	}

	var result []string
	var current strings.Builder
	escaped := false

	for _, r := range s {
		switch {
		case escaped:
			// If we're escaped, always add the character (even if it's the separator)
			current.WriteRune(r)
			escaped = false
		case r == '\\':
			// Start escape sequence
			escaped = true
		case r == sep:
			// Found unescaped separator
			result = append(result, current.String())
			current.Reset()
		default:
			// Regular character
			current.WriteRune(r)
		}
	}

	// Don't forget the last segment
	result = append(result, current.String())

	return result
}
