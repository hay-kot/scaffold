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
func Inject(r io.Reader, data string, at string) ([]byte, error) {
	bldr := strings.Builder{}
	newline := func(s string) {
		bldr.WriteString(s)
		bldr.WriteString("\n")
	}

	found := false

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, at) {
			// Found the line, insert the data
			// before this line
			indent := indentation(line)

			lines := strings.Split(data, "\n")

			for _, l := range lines {
				if l == "" {
					continue
				}

				newline(indent + l)
			}

			found = true
		}

		newline(line)
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	if !found {
		return nil, ErrInjectMarkerNotFound
	}

	return []byte(bldr.String()), nil
}
