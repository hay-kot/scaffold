package scaffold

import "testing"

func TestparseDefaultString(t *testing.T) {
	type tcase struct {
		name     string
		inputs   []any
		expected string
	}

	cases := []tcase{
		{
			name:     "Empty inputs",
			inputs:   []any{},
			expected: "",
		},
		{
			name:     "Single string input",
			inputs:   []any{"test"},
			expected: "test",
		},
		{
			name:     "Multiple inputs, first is string",
			inputs:   []any{nil, "test", true},
			expected: "test",
		},
		{
			name:     "Multiple inputs, no string",
			inputs:   []any{42, false},
			expected: "",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseDefaultString(tc.inputs...)
			if actual != tc.expected {
				t.Errorf("parseDefaultString(%v) = %v; expected %v", tc.inputs, actual, tc.expected)
			}
		})
	}
}

func TestparseDefaultStrings(t *testing.T) {
	type tcase struct {
		name     string
		inputs   []any
		expected []string
	}

	cases := []tcase{
		{
			name:     "No input",
			inputs:   []any{},
			expected: nil,
		},
		{
			name:     "Single string input",
			inputs:   []any{"test"},
			expected: []string{"test"},
		},
		{
			name:     "Multiple inputs with strings",
			inputs:   []any{42, "test1", true, "test2"},
			expected: []string{"test1", "test2"},
		},
		{
			name:     "Multiple inputs without strings",
			inputs:   []any{42, false},
			expected: nil,
		},
		{
			name:     "Nested slice input with strings",
			inputs:   []any{[]any{"test1", "test2"}},
			expected: []string{"test1", "test2"},
		},
		{
			name:     "Nested slice input without strings",
			inputs:   []any{[]any{42, false}},
			expected: nil,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseDefaultStrings(tc.inputs...)
			if len(actual) != len(tc.expected) {
				t.Errorf("parseDefaultStrings(%v) = %v; expected %v", tc.inputs, actual, tc.expected)
			} else {
				for i := 0; i < len(actual); i++ {
					if actual[i] != tc.expected[i] {
						t.Errorf("parseDefaultStrings(%v) = %v; expected %v", tc.inputs, actual, tc.expected)
					}
				}
			}
		})
	}
}

func TestparseDefaultBool(t *testing.T) {
	type tcase struct {
		name     string
		inputs   []any
		expected bool
	}

	cases := []tcase{
		{
			name:     "No input",
			inputs:   []any{},
			expected: false,
		},
		{
			name:     "Single bool input",
			inputs:   []any{true},
			expected: true,
		},
		{
			name:     "Multiple inputs with bool",
			inputs:   []any{"test", false, true},
			expected: false,
		},
		{
			name:     "Multiple inputs without bool",
			inputs:   []any{"test", 42},
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := parseDefaultBool(tc.inputs...)
			if actual != tc.expected {
				t.Errorf("parseDefaultBool(%v) = %v; expected %v", tc.inputs, actual, tc.expected)
			}
		})
	}
}
