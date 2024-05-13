package huhext

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

type LoopedInput struct {
	input *huh.Input
	desc  string

	values   []string
	inputstr *string
}

func NewLoopedInput() *LoopedInput {
	return &LoopedInput{
		input: huh.NewInput(),
	}
}

func (l *LoopedInput) Description(v string) *LoopedInput {
	l.desc = v
	l.input.Description(l.description())
	return l
}

func (l *LoopedInput) Key(v string) *LoopedInput {
	l.input = l.input.Key(v)
	return l
}

func (l *LoopedInput) Title(v string) *LoopedInput {
	l.input = l.input.Title(v)
	return l
}

func (l *LoopedInput) Value(v []string) *LoopedInput {
	l.values = v
	l.input.Description(l.description())
	return l
}

// Blur implements huh.Field.
func (l *LoopedInput) Blur() tea.Cmd {
	return l.input.Blur()
}

// Error implements huh.Field.
func (l *LoopedInput) Error() error {
	return l.input.Error()
}

// Focus implements huh.Field.
func (l *LoopedInput) Focus() tea.Cmd {
	return l.input.Focus()
}

// GetKey implements huh.Field.
func (l *LoopedInput) GetKey() string {
	return l.input.GetKey()
}

// GetValue implements huh.Field.
func (l *LoopedInput) GetValue() any {
	lastval := l.input.GetValue().(string)

	if lastval != "" {
		l.values = append(l.values, lastval)
	}

	return l.values
}

// Init implements huh.Field.
func (l *LoopedInput) Init() tea.Cmd {
	return l.input.Init()
}

// KeyBinds implements huh.Field.
func (l *LoopedInput) KeyBinds() []key.Binding {
	custombinds := []key.Binding{
		key.NewBinding(
			key.WithKeys("up"),
			key.WithHelp("up", "prev value"),
		),
	}

	return append(l.input.KeyBinds(), custombinds...)
}

// Run implements huh.Field.
func (l *LoopedInput) Run() error {
	return l.input.Run()
}

// Skip implements huh.Field.
func (l *LoopedInput) Skip() bool {
	return l.input.Skip()
}

func (l *LoopedInput) description() string {
	if len(l.values) == 0 {
		return l.desc
	}

	return fmt.Sprintf("chosen: %v", strings.Join(l.values, ", "))
}

// Update implements huh.Field.
func (l *LoopedInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if msg, ok := msg.(tea.KeyMsg); ok {
		switch msg.String() {
		case "up": // pop that last value from the slice and update the input
			lastval := ""
			if len(l.values) > 0 {
				lastval = l.values[len(l.values)-1]
				l.values = l.values[:len(l.values)-1]
			}

			l.inputstr = &lastval

			l.input = l.input.Value(l.inputstr)
			return l, l.input.Focus()
		case "enter":
			val := l.input.GetValue().(string)
			if val == "" {
				m, cmd := l.input.Update(msg)
				l.input = m.(*huh.Input)
				return l, cmd
			}

			l.values = append(l.values, val)

			var zero string
			l.inputstr = &zero

			l.input = l.input.Value(l.inputstr)
			return l, l.input.Focus()
		}
	}

	m, cmd := l.input.Update(msg)
	l.input = m.(*huh.Input)
	return l, cmd
}

// View implements huh.Field.
func (l *LoopedInput) View() string {
	l.input = l.input.Description(l.description())
	return l.input.View()
}

// WithAccessible implements huh.Field.
func (l *LoopedInput) WithAccessible(b bool) huh.Field {
	l.input.WithAccessible(b)
	return l
}

// WithHeight implements huh.Field.
func (l *LoopedInput) WithHeight(h int) huh.Field {
	l.input.WithHeight(h)
	return l
}

// WithKeyMap implements huh.Field.
func (l *LoopedInput) WithKeyMap(mp *huh.KeyMap) huh.Field {
	l.input.WithKeyMap(mp)
	return l
}

// WithPosition implements huh.Field.
func (l *LoopedInput) WithPosition(v huh.FieldPosition) huh.Field {
	l.input.WithPosition(v)
	return l
}

// WithTheme implements huh.Field.
func (l *LoopedInput) WithTheme(v *huh.Theme) huh.Field {
	l.input.WithTheme(v)
	return l
}

// WithWidth implements huh.Field.
func (l *LoopedInput) WithWidth(v int) huh.Field {
	l.input.WithWidth(v)
	return l
}
