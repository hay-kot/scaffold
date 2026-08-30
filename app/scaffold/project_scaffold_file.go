package scaffold

import (
	"io"
	"slices"

	"gopkg.in/yaml.v3"
)

type ProjectScaffoldFile struct {
	Metadata   Metadata                  `yaml:"metadata"`
	Skip       []string                  `yaml:"skip"`
	Questions  []Question                `yaml:"questions"`
	Rewrites   []Rewrite                 `yaml:"rewrites"`
	Computed   map[string]string         `yaml:"computed"`
	Messages   Messages                  `yaml:"messages"`
	Inject     []Injectable              `yaml:"inject"`
	Features   []Feature                 `yaml:"features"`
	Presets    map[string]map[string]any `yaml:"presets"`
	Delimiters []Delimiters              `yaml:"delimiters"`
	Each       []EachConfig              `yaml:"each"`

	// computedOrder records the order the computed keys appear in scaffold.yaml.
	// Computed stays a plain map, so the order rides beside it rather than in it.
	// ReadScaffoldFile populates this; a caller that decodes the struct itself
	// leaves it empty and ComputedOrder falls back to sorted keys.
	computedOrder []string
}

// computedKeyOrder reads the keys of the top level "computed" mapping in the
// order they are written. It returns nil for any shape it does not recognise,
// which leaves ComputedOrder to fall back to sorted keys.
func computedKeyOrder(doc *yaml.Node) []string {
	if doc.Kind == yaml.DocumentNode && len(doc.Content) == 1 {
		doc = doc.Content[0]
	}

	if doc.Kind != yaml.MappingNode {
		return nil
	}

	for i := 0; i+1 < len(doc.Content); i += 2 {
		if doc.Content[i].Value != "computed" {
			continue
		}

		node := doc.Content[i+1]
		if node.Kind != yaml.MappingNode {
			return nil
		}

		keys := make([]string, 0, len(node.Content)/2)
		for j := 0; j+1 < len(node.Content); j += 2 {
			keys = append(keys, node.Content[j].Value)
		}

		return keys
	}

	return nil
}

// ComputedOrder returns the computed keys in the order they must resolve, so an
// entry can reference the entries above it. Keys read from scaffold.yaml keep
// their declaration order. Keys added programmatically have no declared order,
// so they follow in sorted order to keep a render reproducible.
func (p *ProjectScaffoldFile) ComputedOrder() []string {
	if len(p.Computed) == 0 {
		return nil
	}

	out := make([]string, 0, len(p.Computed))
	seen := make(map[string]bool, len(p.Computed))

	for _, key := range p.computedOrder {
		if _, ok := p.Computed[key]; ok && !seen[key] {
			out = append(out, key)
			seen[key] = true
		}
	}

	rest := make([]string, 0, len(p.Computed)-len(out))
	for key := range p.Computed {
		if !seen[key] {
			rest = append(rest, key)
		}
	}

	slices.Sort(rest)

	return append(out, rest...)
}

// EachConfig declares a variable for multi-file expansion. It supports both
// a string shorthand ("services") and an object form ({var: "models", as: "..."}).
type EachConfig struct {
	Var string
	As  string
}

func (e *EachConfig) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind == yaml.ScalarNode {
		e.Var = value.Value
		return nil
	}

	var obj struct {
		Var string `yaml:"var"`
		As  string `yaml:"as"`
	}
	if err := value.Decode(&obj); err != nil {
		return err
	}
	e.Var = obj.Var
	e.As = obj.As
	return nil
}

type Delimiters struct {
	Glob  string `yaml:"glob"`
	Left  string `yaml:"left"`
	Right string `yaml:"right"`
}

func ReadScaffoldFile(reader io.Reader) (*ProjectScaffoldFile, error) {
	// Decode through a node rather than straight into the struct. A map cannot
	// hold the order its keys were written in, and computed variables resolve in
	// declaration order, so the order is read off the node and kept beside the
	// map. Decoding the node into the struct keeps the yaml error messages
	// pointing at ProjectScaffoldFile.
	var doc yaml.Node

	err := yaml.NewDecoder(reader).Decode(&doc)
	if err != nil {
		return nil, err
	}

	var out ProjectScaffoldFile
	if err := doc.Decode(&out); err != nil {
		return nil, err
	}

	out.computedOrder = computedKeyOrder(&doc)

	return &out, nil
}

type Messages struct {
	Pre  string `yaml:"pre"`
	Post string `yaml:"post"`
}

type Rewrite struct {
	From string `yaml:"from"`
	To   string `yaml:"to"`
}

type Mode string

const (
	Before Mode = "before"
	After  Mode = "after"
)

type Injectable struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	At       string `yaml:"at"`
	Mode     Mode   `yaml:"mode"`
	Template string `yaml:"template"`
}

type Feature struct {
	Value string   `yaml:"value"`
	Globs []string `yaml:"globs"`
}
