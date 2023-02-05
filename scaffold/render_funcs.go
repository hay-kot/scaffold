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

type ProjectRenderOptions struct {
	OutDirectory string
}

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
		log.Debug().Str("path", node.GetTemplatePath()).Msg("Rendering node")

		if node.folder {
			log.Debug().Str("path", node.GetTemplatePath()).Msg("Creating folder")

			path := node.GetTemplatePath()
			path, err := e.TmplString(path, iVars)
			if err != nil {
				return err
			}

			path = filepath.Join(projectPath, path)

			log.Debug().Str("path", path).Msg("Creating folder")
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
			return fmt.Errorf("failed to render node: %w", err)
		}

		outpath := filepath.Join(projectPath, node.outpath)

		// Write to file
		log.Debug().Str("path", outpath).Msg("Writing to file")
		err = os.WriteFile(outpath, node.outContent, 0644)
		if err != nil {
			return fmt.Errorf("failed to write to file: %w", err)
		}

		close()
	}

	return nil
}
