// Package scaffold provides a simple way to scaffold a project from a template
package scaffold

import (
	"fmt"
	"io"
	"io/fs"
	"maps"
	"strconv"
	"strings"

	"charm.land/huh/v2"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

const (
	HooksDir            = "hooks"
	PostScaffoldScripts = "post_scaffold"
	TemplateDirName     = "templates"
)

var projectNames = [...]string{
	"{{ .Project }}",
	"{{ .ProjectSlug }}",
	"{{ .ProjectSnake }}",
	"{{ .ProjectKebab }}",
	"{{ .ProjectCamel }}",
	"{{ .ProjectPascal }}",
}

var templateNames = [...]string{
	TemplateDirName,
}

// Project structure hold the project templates file system and configuration for
// rendering the project.
type Project struct {
	RootFS       rwfs.ReadFS
	NameTemplate string
	Name         string
	Conf         *ProjectScaffoldFile
	Options      Options
}

func readFirst(fsys fs.FS, names ...string) (fs.File, error) {
	for _, name := range names {
		file, err := fsys.Open(name)
		if err == nil {
			return file, nil
		}
	}

	return nil, fmt.Errorf("file not found")
}

func LoadProject(fileSys fs.FS, opts Options) (*Project, error) {
	p := &Project{
		RootFS:  fileSys,
		Options: opts,
	}

	var err error

	scaffoldFile, err := readFirst(p.RootFS, "scaffold.yaml", "scaffold.yml")
	if err != nil {
		return nil, fmt.Errorf("scaffold.{yml,yaml} does not exist")
	}

	p.NameTemplate, err = RootDir(p.RootFS)
	if err != nil {
		return nil, err
	}

	scaffold, err := ReadScaffoldFile(scaffoldFile)
	if err != nil {
		return nil, err
	}

	p.Conf = scaffold

	return p, nil
}

// RootDir returns the name of the scaffold's root template directory: either
// "templates" or one of the "{{ .Project }}" name variants. It returns an error
// when none exist, or when more than one exists, because only one of them can
// be rendered and the others would be silently ignored.
func RootDir(fsys fs.FS) (string, error) {
	var found []string
	for _, dir := range append(projectNames[:], templateNames[:]...) {
		info, err := fs.Stat(fsys, dir)
		if err == nil && info.IsDir() {
			found = append(found, dir)
		}
	}

	switch len(found) {
	case 0:
		return "", fmt.Errorf("{{ .Project }} or templates directory does not exist")
	case 1:
		return found[0], nil
	default:
		return "", fmt.Errorf("scaffold must have exactly one root directory, found %d: %s", len(found), strings.Join(quoteAll(found), ", "))
	}
}

func quoteAll(ss []string) []string {
	out := make([]string, len(ss))
	for i, s := range ss {
		out[i] = strconv.Quote(s)
	}
	return out
}

// AskQuestions prompts for the project's questions and writes a summary of the
// answers to summary.
func (p *Project) AskQuestions(def map[string]any, e *engine.Engine, theme huh.Theme, summary io.Writer) (map[string]any, error) {
	projectMode := p.NameTemplate != TemplateDirName

	if projectMode {
		name, ok := def["Project"]

		switch ok {
		case false:
			msg := "Project Name"
			decs := "The name your project will be generated with"

			pre := []Question{
				{
					Name:  "Project",
					Group: "default",
					Prompt: AnyPrompt{
						Message:     &msg,
						Description: &decs,
					},
					Required: true,
				},
			}

			p.Conf.Questions = append(pre, p.Conf.Questions...)
		case true:
			nameStr, ok := name.(string)
			if !ok {
				return nil, fmt.Errorf("Project name must be a string")
			}

			p.Name = nameStr
		}
	}

	vars := maps.Clone(def)

	qgroups := QuestionGroupBy(p.Conf.Questions)
	askables := []*Askable{}
	patchvars := func() error {
		for _, askable := range askables {
			err := askable.Hook(vars)
			if err != nil {
				return err
			}
		}

		return nil
	}

	var form *huh.Form
	formgroups := []*huh.Group{}

	for _, qgroup := range qgroups {
		fields := []huh.Field{}

		for _, q := range qgroup {
			question := q.ToAskable(vars[q.Name])
			fields = append(fields, question.Field)
			askables = append(askables, question)
		}

		group := huh.NewGroup(fields...)

		firstq := qgroup[0]
		if firstq.When != "" {
			group.WithHideFunc(func() bool {
				if form == nil {
					return false
				}

				// extract existing properties
				_ = patchvars()

				first := qgroup[0]

				// we check the first question in the group to see if it has a when
				// and if so, we evaluate it and skip the group if it's false
				if first.When != "" {
					result, err := e.TmplString(first.When, vars)
					if err != nil {
						return true
					}

					resultBool, _ := strconv.ParseBool(result)
					if !resultBool {
						return true
					}
				}
				return false
			})
		}

		formgroups = append(formgroups, group)
	}

	form = huh.NewForm(formgroups...).WithTheme(theme)

	err := form.Run()
	if err != nil {
		return nil, err
	}

	// Ensure properts are set on vars
	err = patchvars()
	if err != nil {
		return nil, err
	}

	// Grab the project name from the vars/answers to ensure that
	// it's set.
	if projectMode {
		p.Name = vars["Project"].(string)
	} else {
		p.Name = TemplateDirName
		vars["Project"] = TemplateDirName
	}

	if log.Logger.GetLevel() == zerolog.DebugLevel {
		for k, v := range vars {
			log.Debug().
				Str("key", k).
				Type("type", v).
				Interface("value", v).
				Msg("question")
		}
	}

	for _, askable := range askables {
		_, _ = fmt.Fprint(summary, askable.String())
	}

	_, _ = fmt.Fprint(summary, "\n")

	return vars, nil
}
