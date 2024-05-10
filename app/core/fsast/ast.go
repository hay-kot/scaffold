// Package fsast provides a way to build an abstract syntax tree (AST) of a file system.
package fsast

import (
	"bytes"
	"fmt"
	"io/fs"
	"strings"
)

type NodeType string

const (
	DirNodeType  NodeType = "dir"
	FileNodeType NodeType = "file"
)

type AstNode struct {
	NodeType NodeType
	Path     string
	Content  []byte
	Leafs    []*AstNode
}

func (n *AstNode) string(indent int) string {
	ast := &strings.Builder{}

	indentStr := strings.Repeat("\t", indent)

	for _, child := range n.Leafs {
		prefix := fmt.Sprintf("%s:  (type=%s)", child.Path, child.NodeType)
		if child.NodeType == "file" {
			// ensure all new lines in the file are indented
			content := string(bytes.ReplaceAll(child.Content, []byte("\n"), []byte("\n"+indentStr+"\t")))

			ast.WriteString(indentStr + prefix + "\n" + indentStr + "\t" + content + "\n")
		} else {
			ast.WriteString(indentStr + prefix + "\n" + child.string(indent+1))
		}
	}

	return ast.String()
}

func (n *AstNode) String() string {
	return n.string(0)
}

func New(subFs fs.FS) (*AstNode, error) {
	root := &AstNode{
		NodeType: DirNodeType,
		Path:     "ROOT_NODE",
	}

	return root, Build(subFs, root)
}

func Build(subFs fs.FS, root *AstNode) error {
	files, err := fs.ReadDir(subFs, ".")
	if err != nil {
		return err
	}

	for _, file := range files {
		node := &AstNode{
			Path: file.Name(),
		}

		if file.IsDir() {
			node.NodeType = DirNodeType
			subsubFS, _ := fs.Sub(subFs, file.Name())

			err = Build(subsubFS, node)
			if err != nil {
				return err
			}
		} else {
			node.NodeType = FileNodeType
			readContent, err := fs.ReadFile(subFs, file.Name())
			if err != nil {
				return err
			}
			node.Content = readContent
		}

		root.Leafs = append(root.Leafs, node)
	}

	return nil
}
