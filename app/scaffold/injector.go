package scaffold

import (
	"bufio"
	"bytes"
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
	var buf bytes.Buffer

	// Write a line to the buffer
	writeLine := func(line string) {
		buf.WriteString(line)
		buf.WriteString("\n")
	}

	// Write multiple lines with indentation to the buffer
	writeLines := func(lines []string, indent string) {
		for _, l := range lines {
			if l != "" {
				writeLine(indent + l)
			}
		}
	}

	var (
		inserted      = false
		found         = false
		scanner       = bufio.NewScanner(r)
		linesToInsert = strings.Split(data, "\n")
	)

	// Loop through each line in the reader
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, at) {
			if mode != After {
				writeLines(linesToInsert, indentation(line))
				inserted = true
			}
			found = true
		}

		writeLine(line)

		// If in 'after' mode and the insertion point is found, insert after it
		if mode == After && found && !inserted {
			writeLines(linesToInsert, indentation(line))
			inserted = true
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrInjectMarkerNotFound
	}

	return buf.Bytes(), nil
}
