package main

import (
	"log"
	"os"

	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
)

func main() {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("name").
				Title("Name").
				Description("The name of the project"),
			input(),
			huh.NewInput().
				Key("name").
				Title("Name").
				Description("The name of the project"),
		),
	)

	err := form.Run()
	if err != nil {
		panic(err)
	}

	v := form.Get("names")

	// cast as []string
	names := v.([]string)

	for i, name := range names {
		println(i, name)
	}

	// Do something with the form data
	println("all done")
}

func input() huh.Field {
	debugfile, err := os.Create("debug.log")
	if err != nil {
		panic(err)
	}

	loggers := log.New(debugfile, "", log.LstdFlags)

	input := huh.NewInput().
		Key("names").
		Title("Name").
		Description("The name of the project")

	return &LoopedInput{
		input: input,
		log:   loggers,
	}
}

var _ huh.Field = &LoopedInput{}

type LoopedInput struct {
	input    *huh.Input
	values   []string
	inputstr *string
	log      *log.Logger
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
	return l.input.KeyBinds()
}

// Run implements huh.Field.
func (l *LoopedInput) Run() error {
	return l.input.Run()
}

// Skip implements huh.Field.
func (l *LoopedInput) Skip() bool {
	return l.input.Skip()
}

// Update implements huh.Field.
func (l *LoopedInput) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		l.log.Printf("msg value: %v\n", msg.String())
		if msg.String() == "enter" {
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
	default:
		l.log.Printf("msg type: %T\n", msg)
	}

	m, cmd := l.input.Update(msg)
	l.input = m.(*huh.Input)
	return l, cmd
}

// View implements huh.Field.
func (l *LoopedInput) View() string {
	l.log.Println("View")
	return l.input.View()
}

// WithAccessible implements huh.Field.
func (l *LoopedInput) WithAccessible(b bool) huh.Field {
	l.input = l.input.WithAccessible(b).(*huh.Input)
	return l
}

// WithHeight implements huh.Field.
func (l *LoopedInput) WithHeight(h int) huh.Field {
	l.input = l.input.WithHeight(h).(*huh.Input)
	return l
}

// WithKeyMap implements huh.Field.
func (l *LoopedInput) WithKeyMap(mp *huh.KeyMap) huh.Field {
	l.input = l.input.WithKeyMap(mp).(*huh.Input)
	return l
}

// WithPosition implements huh.Field.
func (l *LoopedInput) WithPosition(v huh.FieldPosition) huh.Field {
	l.input = l.input.WithPosition(v).(*huh.Input)
	return l
}

// WithTheme implements huh.Field.
func (l *LoopedInput) WithTheme(v *huh.Theme) huh.Field {
	l.input = l.input.WithTheme(v).(*huh.Input)
	return l
}

// WithWidth implements huh.Field.
func (l *LoopedInput) WithWidth(v int) huh.Field {
	l.input = l.input.WithWidth(v).(*huh.Input)
	return l
}
