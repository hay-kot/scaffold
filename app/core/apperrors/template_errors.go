package apperrors

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/hay-kot/scaffold/internal/styles"
)

// TemplateError represents an error that occurred during template processing
type TemplateError struct {
	FilePath     string
	LineNumber   int
	ColumnNumber int
	Original     error
	Context      []string // Lines of code around the error
	Delimiters   DelimiterInfo
}

type DelimiterInfo struct {
	Left  string
	Right string
}

func (e *TemplateError) Error() string {
	var b strings.Builder
	b.WriteString(e.Original.Error())
	if e.FilePath != "" {
		b.WriteString(fmt.Sprintf(" in %s", e.FilePath))
		if e.LineNumber > 0 {
			b.WriteString(fmt.Sprintf(":%d", e.LineNumber))
			if e.ColumnNumber > 0 {
				b.WriteString(fmt.Sprintf(":%d", e.ColumnNumber))
			}
		}
	}
	return b.String()
}

// ConsoleOutput provides a rich formatted error message for the console
// This implements the printer.ConsoleOutput interface
func (e *TemplateError) ConsoleOutput() string {
	var b strings.Builder

	// Error header
	errorHeaderStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ColorError)).
		Bold(true)

	b.WriteString(errorHeaderStyle.Render(fmt.Sprintf("%s Template Error", styles.Cross)))
	b.WriteString("\n")

	// Create the error box with left border only
	boxStyle := lipgloss.NewStyle().
		BorderLeft(true).
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("243")).
		PaddingLeft(2)

	var boxContent strings.Builder

	// Error message
	errorMsgStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Bold(true)
	boxContent.WriteString(errorMsgStyle.Render(e.getCleanErrorMessage()))
	boxContent.WriteString("\n")

	// File location
	if e.FilePath != "" {
		fileStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("69"))
		location := filepath.Base(e.FilePath)
		if e.LineNumber > 0 {
			location = fmt.Sprintf("%s:%d", location, e.LineNumber)
		}

		boxContent.WriteString(fileStyle.Render(location))

		// Show relative path if different from basename
		if e.FilePath != filepath.Base(e.FilePath) {
			pathStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
			boxContent.WriteString(pathStyle.Render(fmt.Sprintf(" (%s)", e.FilePath)))
		}

		// Code context if available
		if len(e.Context) > 0 && e.LineNumber > 0 {
			boxContent.WriteString("\n")

			lineNumStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
			errorLineStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color(styles.ColorError)).
				Bold(true)
			codeStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

			// Calculate starting line number based on context
			startLine := e.LineNumber
			if len(e.Context) == 3 {
				// We have before, current, after
				startLine = e.LineNumber - 1
			} else if len(e.Context) == 2 && e.LineNumber > 1 {
				// We have before and current (no after line)
				startLine = e.LineNumber - 1
			}
			// else we only have the error line or error line + after

			for i, line := range e.Context {
				lineNum := startLine + i
				isErrorLine := lineNum == e.LineNumber

				if isErrorLine {
					// Error line with arrow
					boxContent.WriteString(errorLineStyle.Render(fmt.Sprintf("→ %3d │ ", lineNum)))
					boxContent.WriteString(errorLineStyle.Render(line))
				} else {
					// Context line
					boxContent.WriteString(lineNumStyle.Render(fmt.Sprintf("  %3d │ ", lineNum)))
					boxContent.WriteString(codeStyle.Render(line))
				}

				if i < len(e.Context)-1 {
					boxContent.WriteString("\n")
				}
			}
		}

		// Show delimiters if available
		if e.Delimiters.Left != "" && e.Delimiters.Right != "" {
			boxContent.WriteString("\n\n")
			delimStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("243"))
			boxContent.WriteString(delimStyle.Render(fmt.Sprintf("delimiters: %s ... %s", e.Delimiters.Left, e.Delimiters.Right)))
		}
	}

	b.WriteString(boxStyle.Render(boxContent.String()))
	b.WriteString("\n\n")

	return b.String()
}

// getCleanErrorMessage cleans up Go template error messages to be more readable
func (e *TemplateError) getCleanErrorMessage() string {
	errMsg := e.Original.Error()

	// Remove the "template: scaffold:123:" prefix if present
	if strings.HasPrefix(errMsg, "template:") {
		parts := strings.SplitN(errMsg, ": ", 3)
		if len(parts) >= 3 {
			// The actual error message is usually the last part
			errMsg = parts[len(parts)-1]
		} else if len(parts) == 2 {
			errMsg = parts[1]
		}
	}

	// Clean up common template error messages to be more readable
	replacements := map[string]string{
		"function \"":     "function '",
		"\" not defined":  "' is not defined",
		"unexpected ":     "unexpected token ",
		"unclosed action": "unclosed template action (missing closing delimiter)",
		"bad character":   "invalid character",
		"expected":        "expected",
	}

	for old, new := range replacements {
		errMsg = strings.ReplaceAll(errMsg, old, new)
	}

	return errMsg
}

// WrapTemplateError wraps an error with template-specific context
func WrapTemplateError(err error, filePath string) *TemplateError {
	if err == nil {
		return nil
	}

	// If it's already a TemplateError, just update the file path if needed
	var te *TemplateError
	if errors.As(err, &te) {
		if filePath != "" && te.FilePath == "" {
			te.FilePath = filePath
		}
		return te
	}

	te = &TemplateError{
		FilePath: filePath,
		Original: err,
	}

	// Try to extract line and column numbers from the error message
	te.extractLocationFromError()

	return te
}

// extractLocationFromError attempts to extract line and column numbers from the error message
func (e *TemplateError) extractLocationFromError() {
	errMsg := e.Original.Error()

	// Look for patterns like "template:name:123:45" or "template:name:123"
	// Common in Go template errors
	if strings.Contains(errMsg, "template:") {
		parts := strings.Split(errMsg, ":")
		for i := 1; i < len(parts)-1; i++ {
			part := strings.TrimSpace(parts[i])
			if lineNum, err := strconv.Atoi(part); err == nil && lineNum > 0 {
				e.LineNumber = lineNum
				// Try to get column number from the next part
				if i+1 < len(parts)-1 {
					nextPart := strings.TrimSpace(parts[i+1])
					if colNum, err := strconv.Atoi(nextPart); err == nil && colNum > 0 {
						e.ColumnNumber = colNum
					}
				}
				break
			}
		}
	}
}

// WithContext adds code context lines to the error
func (e *TemplateError) WithContext(lines []string) *TemplateError {
	e.Context = lines
	return e
}

// WithDelimiters adds delimiter information to the error
func (e *TemplateError) WithDelimiters(left, right string) *TemplateError {
	e.Delimiters = DelimiterInfo{
		Left:  left,
		Right: right,
	}
	return e
}
