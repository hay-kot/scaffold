package scaffold

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/huandu/xstrings"
	"github.com/rs/zerolog/log"
)

type RenderableNode interface {
	GetTemplatePath() string
	SetOutPath(string)
	io.ReadWriter
}

// RenderNode renders a node and sets the outpath on the node while
// processing it as a template.
func RenderNode[T RenderableNode](s *Engine, n T, vars any) error {
	outpath, err := s.TmplString(n.GetTemplatePath(), vars)
	if err != nil {
		return err
	}

	n.SetOutPath(outpath)

	err = RenderReadWriter(s, n, vars)
	if err != nil {
		return err
	}

	return nil
}

// RenderReadWriter renders a io.ReadWriter as a template.
//
//   - The io.Reader is expected to be a template.
//   - The io.Writer is expected to be a destination for the rendered template. (e.g. a file)
func RenderReadWriter[T io.ReadWriter](s *Engine, rw T, vars any) error {
	tmpl, err := s.TmplFactory(rw)
	if err != nil {
		return err
	}

	err = s.RenderTemplate(rw, tmpl, vars)
	if err != nil {
		return err
	}

	return nil
}

// ProjectRenderOptions are options for rendering a project.
type ProjectRenderOptions struct {
	// OutDirectory is the directory to render the project to.
	//
	// Example: /home/user/projects
	//
	//  | /home/user/projects/my-project
	//	|___ <- New project rendered here
	OutDirectory string
}

// RenderProject renders an entire Project to the given out directory.
//
// It will also transforms the vars variable to be used in the templates
// into a nested map of the following structure:
//
//   - Project: The project name
//   - ProjectSnake: The project name in snake case
//   - ProjectKebab: The project name in kebab case
//   - ProjectCamel: The project name in camel case
//   - Scaffold: The vars variable passed to this function
func RenderProject(e *Engine, p *Project, vars any, opts ProjectRenderOptions) error {
	iVars := Vars{
		"Project":      p.qProject,
		"ProjectSnake": xstrings.ToSnakeCase(p.qProject),
		"ProjectKebab": xstrings.ToKebabCase(p.qProject),
		"ProjectCamel": xstrings.ToCamelCase(p.qProject),
		"Scaffold":     vars,
	}

	nodes := p.Tree.Flatten()
	log.Debug().Int("nodes", len(nodes)).Msg("Rendering project")

	projectName, err := e.TmplString(p.ProjectNameTemplate, iVars)
	if err != nil {
		return err
	}

	projectPath := filepath.Join(opts.OutDirectory, projectName)
	err = os.MkdirAll(projectPath, 0755)
	if err != nil && !os.IsExist(err) {
		return err
	}

	for _, node := range nodes {
		if node.folder {
			path := node.GetTemplatePath()
			path, err := e.TmplString(path, iVars)
			if err != nil {
				return err
			}

			path = filepath.Join(projectPath, path)

			err = os.Mkdir(path, 0755)
			if err != nil && !os.IsExist(err) {
				return err
			}

			continue
		}

		close, err := node.Open()
		if err != nil {
			return fmt.Errorf("failed to open node: %w", err)
		}

		err = RenderNode(e, node, iVars)
		if err != nil {
			close()
			return fmt.Errorf("failed to render node: %w", err)
		}

		outpath := filepath.Join(projectPath, node.outpath)

		err = os.WriteFile(outpath, node.outContent, 0644)
		if err != nil {
			close()
			return fmt.Errorf("failed to write to file: %w", err)
		}

		close()
	}

	return nil
}
