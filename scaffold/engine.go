// Package scaffold provides a simple templating engine for scaffolding projects.
package scaffold

import (
	"fmt"
	"io"
	"strings"
	"text/template"

	"github.com/Masterminds/sprig/v3"
	"github.com/rs/zerolog/log"
)

type Vars map[string]interface{}

type Engine struct {
	baseTemplate *template.Template
}

func NewEngine() *Engine {
	return &Engine{
		baseTemplate: template.New("scaffold").Funcs(sprig.FuncMap()),
	}
}

func isTemplate(s string) bool {
	return strings.Contains(s, "{{")
}

func (e *Engine) TmplString(str string, vars any) (string, error) {
	if !isTemplate(str) {
		return str, nil
	}

	tmpl, err := e.baseTemplate.Parse(str)
	if err != nil {
		log.Err(err).Msg("failed to parse template")
		return "", err
	}

	out := &strings.Builder{}

	err = e.RenderTemplate(out, tmpl, vars)
	if err != nil {
		log.Err(err).Msg("failed to render template")
		return "", err
	}

	return out.String(), nil
}

func (e *Engine) TmplFactory(reader io.Reader) (*template.Template, error) {
	if reader == nil {
		return nil, fmt.Errorf("reader is nil")
	}

	out, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	return e.baseTemplate.Parse(string(out))
}

func (e *Engine) RenderTemplate(w io.Writer, tmpl *template.Template, vars any) error {
	err := tmpl.Execute(w, vars)

	if err != nil {
		return err
	}

	return nil
}
