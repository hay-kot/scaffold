// Package printer provides a printer abstraction for human readable output.
package printer

import (
	"io"
	"strings"

	"github.com/hay-kot/scaffold/internal/styles"
)

type ConsoleOutput interface {
	ConsoleOutput() string
}

type Printer struct {
	writer io.Writer
	base   styles.RenderFunc
	light  styles.RenderFunc
}

func New(writer io.Writer) *Printer {
	base, light := styles.ThemeColorsScaffold.Compile()
	return &Printer{
		writer: writer,
		base:   base,
		light:  light,
	}
}

func (c *Printer) WithBase(style styles.RenderFunc) *Printer {
	c.base = style
	return c
}

func (c *Printer) WithLight(style styles.RenderFunc) *Printer {
	c.light = style
	return c
}

func (c *Printer) write(s string) {
	_, _ = c.writer.Write([]byte(s))
}

func (c *Printer) Print(s string) {
	c.write(s)
}

func (c *Printer) Println(s string) {
	c.write(s)
	c.write("\n")
}

// UnknownError printer an error message for an unknown or unexpected error.
// This is used when an error in the system was unexpected, and the error output
// should be displayed to the user.
//
// If the error implements the ConsoleOutput interface, the ConsoleOutput method
// will be called to get the error output.
func (c *Printer) UnknownError(title string, err error) {
	bldr := &strings.Builder{}

	bldr.WriteString(styles.Error("An unexpected error occurred"))
	bldr.WriteString("\n")

	consoleErr, ok := err.(ConsoleOutput)
	if ok {
		bldr.WriteString(consoleErr.ConsoleOutput())
		c.write(bldr.String())
		return
	}

	bldr.WriteString(styles.Padding(styles.Bold(err.Error())))
	bldr.WriteString("\n")

	c.write(bldr.String())
}

func (c *Printer) Title(title string) {
	c.write(styles.Bold(title))
	c.write("\n")
}

type StatusListItem struct {
	Ok     bool
	Status string
}

// StatusList prints a list of status items with a title.
//
// Example:
//
//	Some Title
//	 ✔ Status 1
//	 ✘ Status 2
//	 ✔ Status 3
func (c *Printer) StatusList(title string, items []StatusListItem) {
	bldr := strings.Builder{}

	bldr.WriteString(styles.Padding(styles.Bold(c.base(title))))
	bldr.WriteString("\n")

	for _, item := range items {
		bldr.WriteString("  ")
		if item.Ok {
			bldr.WriteString(
				styles.Success(styles.Check),
			)
		} else {
			bldr.WriteString(styles.Error(styles.Cross))
		}

		bldr.WriteString(" ")
		bldr.WriteString(item.Status)
		bldr.WriteString("\n")
	}

	c.write(bldr.String())
}

// List prints a list of items with a title.
//
//	Example:
//
//	Some Title
//	  - Item 1
//	  - Item 2
//	  - Item 3
func (c *Printer) List(title string, items []string) {
	bldr := strings.Builder{}

	bldr.WriteString(styles.Padding(styles.Bold(c.base(title))))
	bldr.WriteString("\n")

	for _, item := range items {
		bldr.WriteString("   ")
		bldr.WriteString(styles.Dot)
		bldr.WriteString(" ")
		bldr.WriteString(item)
		bldr.WriteString("\n")
	}

	c.write(bldr.String())
}

func (c *Printer) LineBreak() {
	c.write("\n")
}

type KeyValueError struct {
	Key     string
	Message string
}

func (c *Printer) KeyValueValidationError(title string, errors []KeyValueError) {
	bldr := strings.Builder{}

	bldr.WriteString(styles.Error(styles.Bold(title)))
	bldr.WriteString("\n")

	for _, err := range errors {
		bldr.WriteString("  ")
		bldr.WriteString(styles.Error(styles.Cross))
		bldr.WriteString(" ")
		bldr.WriteString(err.Key)
		bldr.WriteString(": ")
		bldr.WriteString(styles.Subtle(err.Message))
		bldr.WriteString("\n")
	}

	c.write(bldr.String())
}
