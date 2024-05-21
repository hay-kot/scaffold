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
}

func New(writer io.Writer) *Printer {
	return &Printer{
		writer: writer,
	}
}

func (c *Printer) write(s string) {
	_, _ = c.writer.Write([]byte(s))
}

func (c *Printer) UnknownError(title string, err error) {
	bldr := &strings.Builder{}

	bldr.WriteString(styles.Error("An unexpected error occurred"))
	bldr.WriteString("\n\n")

	consoleErr, ok := err.(ConsoleOutput)
	if ok {
		bldr.WriteString(consoleErr.ConsoleOutput())
		c.write(bldr.String())
		return
	}

	bldr.WriteString(styles.Padding("Error"))
	bldr.WriteString("\n  '")
	bldr.WriteString(err.Error())
	bldr.WriteString("'\n")

	c.write(bldr.String())
}

func (c *Printer) Title(title string) {
	c.write(styles.Bold(title))
	c.write("\n")
}

func (c *Printer) Print(s string) {
	c.write(s)
}

func (c *Printer) Println(s string) {
	c.write(s)
	c.write("\n")
}

type StatusListItem struct {
	StatusOk bool
	Status   string
}

func (c *Printer) StatusList(title string, items []StatusListItem) {
	bldr := strings.Builder{}

	bldr.WriteString(styles.Padding(styles.Bold(title)))
	bldr.WriteString("\n")

	for _, item := range items {
		bldr.WriteString(" ")
		if item.StatusOk {
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

func (c *Printer) LineBreak() {
	c.write("\n")
}
