// Package engine provides a simple templating engine for scaffolding projects.
package engine

import (
	"fmt"
	"io"
	"strings"
	"text/template"
	"unicode"

	"github.com/gertd/go-pluralize"
	"github.com/go-sprout/sprout/sprigin"
	"github.com/rs/zerolog/log"
)

// ErrTemplateIsEmpty is returned when a provided reader is empty.
var ErrTemplateIsEmpty = fmt.Errorf("template is empty")

type Vars map[string]any

type opts struct {
	delimLeft  string
	delimRight string
}

func WithDelims(left string, right string) func(*opts) {
	return func(o *opts) {
		o.delimLeft = left
		o.delimRight = right
	}
}

type Engine struct {
	fm template.FuncMap
}

func New() *Engine {
	fm := sprigin.FuncMap()

	client := pluralize.NewClient()

	fm["wraptmpl"] = wraptmpl
	fm["isPlural"] = client.IsPlural
	fm["isSingular"] = client.IsSingular
	fm["toSingular"] = client.Singular
	fm["toPlural"] = client.Plural

	return &Engine{
		fm: fm,
	}
}

func (e *Engine) parse(tmpl string, opt opts) (*template.Template, error) {
	return template.New("scaffold").
		Funcs(e.fm).
		Delims(opt.delimLeft, opt.delimRight).
		Parse(tmpl)
}

func isTemplate(s string) bool {
	return strings.Contains(s, "{{")
}

func (e *Engine) TmplString(str string, vars any) (string, error) {
	if !isTemplate(str) {
		return str, nil
	}

	opt := opts{
		delimLeft:  "{{",
		delimRight: "}}",
	}

	tmpl, err := e.parse(str, opt)
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
func (e *Engine) Factory(reader io.Reader, opfns ...func(*opts)) (*template.Template, error) {
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

	opt := opts{
		delimLeft:  "{{",
		delimRight: "}}",
	}

	for _, fn := range opfns {
		fn(&opt)
	}

	return e.parse(string(out), opt)
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
