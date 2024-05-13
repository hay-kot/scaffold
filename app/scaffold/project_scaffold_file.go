package scaffold

import (
	"io"

	"gopkg.in/yaml.v3"
)

type ProjectScaffoldFile struct {
	Skip      []string                  `yaml:"skip"`
	Questions []Question                `yaml:"questions"`
	Rewrites  []Rewrite                 `yaml:"rewrites"`
	Computed  map[string]string         `yaml:"computed"`
	Messages  Messages                  `yaml:"messages"`
	Inject    []Injectable              `yaml:"inject"`
	Features  []Feature                 `yaml:"features"`
	Presets   map[string]map[string]any `yaml:"presets"`
}

func ReadScaffoldFile(reader io.Reader) (*ProjectScaffoldFile, error) {
	var out ProjectScaffoldFile

	err := yaml.NewDecoder(reader).Decode(&out)
	if err != nil {
		return nil, err
	}

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
	After Mode = "after"
)

type Injectable struct {
	Name     string `yaml:"name"`
	Path     string `yaml:"path"`
	At       string `yaml:"at"`
	Mode     Mode `yaml:"mode"`
	Template string `yaml:"template"`
}

type Feature struct {
	Value string   `yaml:"value"`
	Globs []string `yaml:"globs"`
}
