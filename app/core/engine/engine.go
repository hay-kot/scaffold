// Package engine provides a simple templating engine for scaffolding projects.
package engine

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
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
	fm       template.FuncMap
	partials map[string]*template.Template // Store partial templates
}

func New() *Engine {
	fm := sprigin.FuncMap()

	e := &Engine{
		fm:       fm,
		partials: map[string]*template.Template{},
	}

	// Template Utilities
	fm["wraptmpl"] = wraptmpl

	// Pluralize
	client := pluralize.NewClient()

	fm["isPlural"] = client.IsPlural
	fm["isSingular"] = client.IsSingular
	fm["toSingular"] = client.Singular
	fm["toPlural"] = client.Plural

	// Re-usable Partials
	fm["partial"] = func(name string, data any) (string, error) {
		tmpl, exists := e.partials[name]
		if !exists {
			return "", fmt.Errorf("partial not found: %s", name)
		}

		var buf strings.Builder
		err := tmpl.Execute(&buf, data)
		if err != nil {
			return "", err
		}

		return buf.String(), nil
	}

	return e
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

// RegisterPartialsFS walks a [fs.FS] and registers all files as partials for the engine.
// It will use a relative path from the specified directory for each partial within the fs.
// If dirname is provided, that prefix will be removed from partial names.
//
// Example with dirname = "templates":
//
//	├── templates/
//	│   ├── common/
//	│   │    ├── snippet1.tmpl
//	│   │    └── snippet2.tmpl
//	│   ├── header.tmpl
//	│   └── footer.tmpl
//
// This structure will register the following partials
// - common/snippet1
// - common/snippet2
// - header
// - footer
func (e *Engine) RegisterPartialsFS(rfs fs.FS, dirname string) error {
	return fs.WalkDir(rfs, dirname, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("failed to walk partials directory: %w", err)
		}

		if d.IsDir() {
			return nil
		}

		content, err := fs.ReadFile(rfs, path)
		if err != nil {
			return fmt.Errorf("failed to read partial file %s: %w", path, err)
		}

		// Get the partial name by removing the file extension and prefix
		name := path

		// Remove file extension
		if ext := filepath.Ext(name); ext != "" {
			name = strings.TrimSuffix(name, ext)
		}

		name = filepath.ToSlash(name)

		return e.RegisterPartial(name, string(content))
	})
}

func (e *Engine) RegisterPartial(name string, content string, opfns ...func(*opts)) error {
	if !IsValidIdentifier(name) {
		return fmt.Errorf("invalid partial name: %s", name)
	}

	opt := opts{
		delimLeft:  "{{",
		delimRight: "}}",
	}

	for _, fn := range opfns {
		fn(&opt)
	}

	tmpl, err := e.parse(content, opt)
	if err != nil {
		return fmt.Errorf("failed to parse partial template %s: %w", name, err)
	}

	log.Debug().Str("name", name).Msg("registering partial")
	e.partials[name] = tmpl
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
