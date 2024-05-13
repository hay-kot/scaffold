// Package engine provides a simple templating engine for scaffolding projects.
package engine

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/go-sprout/sprout"
	"github.com/rs/zerolog/log"
)

// ErrTemplateIsEmpty is returned when a provided reader is empty.
var ErrTemplateIsEmpty = fmt.Errorf("template is empty")

type Vars map[string]any

type Engine struct {
	fm template.FuncMap
}

func New() *Engine {
	fm := sprout.FuncMap()

	fm["wraptmpl"] = wraptmpl

	return &Engine{
		fm: fm,
	}
}

func (e *Engine) parse(tmpl string) (*template.Template, error) {
	return template.New("scaffold").Funcs(e.fm).Parse(tmpl)
}

func isTemplate(s string) bool {
	return strings.Contains(s, "{{")
}

func (e *Engine) TmplString(str string, vars any) (string, error) {
	if !isTemplate(str) {
		return str, nil
	}

	tmpl, err := e.parse(str)
	if err != nil {
		log.Err(err).Msg("failed to parse template")
		return "", err
	}

	out := &strings.Builder{}

	err = e.Render(out, tmpl, vars)
	if err != nil {
		log.Err(err).Msg("failed to render template")
		return "", err
	}

	return out.String(), nil
}

// Factory returns a new template from the provided reader.
// if the reader is empty, an ErrTemplateIsEmpty is returned.
func (e *Engine) Factory(reader io.Reader) (*template.Template, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader is nil")
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	if len(out) == 0 {
		return nil, ErrTemplateIsEmpty
	}

	return e.parse(string(out))
}

func (e *Engine) Render(w io.Writer, tmpl *template.Template, vars any) error {
	err := tmpl.Execute(w, vars)
	if err != nil {
		return err
	}

	return nil
}

func IsValidIdentifier(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' {
			return false
		}
	}
	return true
}
