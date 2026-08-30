package huhext_test

import (
	"slices"
	"testing"

	tea "charm.land/bubbletea/v2"
	"github.com/hay-kot/scaffold/internal/huhext"
)

func press(text string) tea.KeyPressMsg {
	return tea.KeyPressMsg{Code: rune(text[0]), Text: text}
}

func typeText(t *testing.T, input *huhext.LoopedInput, text string) {
	t.Helper()

	for _, r := range text {
		input.Update(press(string(r)))
	}
}

func values(t *testing.T, input *huhext.LoopedInput) []string {
	t.Helper()

	v, ok := input.GetValue().([]string)
	if !ok {
		t.Fatalf("GetValue returned %T, want []string", input.GetValue())
	}

	return v
}

// TestLoopedInputCollectsValues drives the field the way a form does: enter
// commits the current text and clears the input, up pops the last committed
// value back into it.
func TestLoopedInputCollectsValues(t *testing.T) {
	input := huhext.NewLoopedInput().Title("Langs")
	input.Init()
	input.Focus()

	typeText(t, input, "go")
	input.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	typeText(t, input, "rust")
	input.Update(tea.KeyPressMsg{Code: tea.KeyEnter})

	if got := values(t, input); !slices.Equal(got, []string{"go", "rust"}) {
		t.Fatalf("after two entries: got %v, want [go rust]", got)
	}

	input.Update(tea.KeyPressMsg{Code: tea.KeyUp})

	if got := values(t, input); !slices.Equal(got, []string{"go", "rust"}) {
		t.Fatalf("after up: got %v, want [go rust] with rust back in the input", got)
	}
}
