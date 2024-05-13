package scaffold

import (
	"bufio"
	"errors"
	"io"
	"strings"
)

var ErrInjectMarkerNotFound = errors.New("inject marker not found")

// indentation returns the indentation of the line
func indentation(b string) string {
	var i int
	for i = 0; i < len(b); i++ {
		if b[i] != ' ' && b[i] != '\t' {
			break
		}
	}
	return b[:i]
}

// Inject will read the reader line by line and find the line
// that contains the string "at". It will then insert the data
// before that line.
func Inject(r io.Reader, data string, at string, mode Mode) ([]byte, error) {
	bldr := strings.Builder{}
	newline := func(s string) {
		bldr.WriteString(s)
		bldr.WriteString("\n")
	}

	writelines := func(lines []string, indent string) {
		for _, l := range lines {
			if l == "" {
				continue
			}

			newline(indent + l)
		}
	}

	found := false
	after := mode == After
	inserted := false

	scanner := bufio.NewScanner(r)
	lines := strings.Split(data, "\n")
	var indent string

	for scanner.Scan() {
		line := scanner.Text()

		// Found the line, insert the data
		// default case will be before
		if strings.Contains(line, at) {
			indent = indentation(line)
			if mode != After {
				writelines(lines, indent)
				inserted = true
			}
			found = true
		}

		newline(line)

		// if there is an after mode, insert after!
		if after && found && !inserted {
			println("inserting after")
			writelines(lines, indent)
			inserted = true
		}

	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrInjectMarkerNotFound
	}

	return []byte(bldr.String()), nil
}
