package scaffold

import (
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
)

type TemplateNode struct {
	rootFS fs.FS
	folder bool
	// templatePath to the template file
	templatePath string
	templateFile io.Reader

	outpath    string
	outContent []byte
	Children   []*TemplateNode
}

func (t *TemplateNode) Open() (func(), error) {
	if t.folder {
		return func() {}, nil
	}

	file, err := t.rootFS.Open(t.templatePath)
	if err != nil {
		return func() {}, fmt.Errorf("failed to open template file: %w", err)
	}

	t.templateFile = file
	return func() { _ = file.Close() }, nil
}

func (t *TemplateNode) GetTemplatePath() string {
	return t.templatePath
}

func (t *TemplateNode) SetOutPath(path string) {
	t.outpath = path
}

func (t *TemplateNode) Read(p []byte) (n int, err error) {
	if t.folder {
		return 0, io.EOF
	}

	if t.templateFile == nil {
		panic("template file is nil, did you forget to call Open()?")
	}

	return t.templateFile.Read(p)
}

func (t *TemplateNode) Write(p []byte) (int, error) {
	if t.folder {
		return len(p), nil
	}

	t.outContent = append(t.outContent, p...)
	return len(p), nil
}

func (t *TemplateNode) Flatten() []*TemplateNode {
	nodes := []*TemplateNode{t}

	if t.folder {
		for _, child := range t.Children {
			nodes = append(nodes, child.Flatten()...)
		}
	}

	return nodes
}

func (t *TemplateNode) OutPath(e *Engine, vars Vars) (string, error) {
	return e.TmplString(t.templatePath, vars)
}

func parseTemplateNodeTree(fileSys fs.FS, path string) (*TemplateNode, error) {
	// Check if the path is a folder or a file
	isFolder, err := isDirectory(fileSys, path)
	if err != nil {
		return nil, fmt.Errorf("failed to check if path is a directory: %w", err)
	}

	node := &TemplateNode{
		rootFS:       fileSys,
		folder:       isFolder,
		templatePath: path,
		Children:     []*TemplateNode{},
	}

	// If the path is a folder, parse its children
	if isFolder {
		children, err := readDirectory(fileSys, path)
		if err != nil {
			return nil, err
		}

		for _, child := range children {
			childPath := filepath.Join(path, child)
			childNode, err := parseTemplateNodeTree(fileSys, childPath)
			if err != nil {
				return nil, err
			}
			node.Children = append(node.Children, childNode)
		}
	}

	return node, nil
}

// Helper function to check if a path is a directory
func isDirectory(fileSys fs.FS, path string) (bool, error) {
	fi, err := fs.Stat(fileSys, path)
	if err != nil {
		return false, err
	}
	return fi.IsDir(), nil
}

// Helper function to read the contents of a directory
func readDirectory(fileSys fs.FS, path string) ([]string, error) {
	children, err := fs.ReadDir(fileSys, path)
	if err != nil {
		return nil, err
	}

	names := make([]string, len(children))
	for i, child := range children {
		names[i] = child.Name()
	}

	return names, nil
}
