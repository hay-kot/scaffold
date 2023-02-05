package scaffold

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"strings"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/require"
)

func Test_RenderTree(t *testing.T) {
	tests := []struct {
		name string
		fs   fs.FS
	}{
		{
			name: "basic",
			fs:   DynamicFiles(),
		},
		{
			name: "nested",
			fs:   NestedFiles(),
		},
	}

	vars := Vars{
		"Name":  "Your Name1",
		"Name2": "Your Name2",
	}

	snapshot := cupaloy.New(
		cupaloy.SnapshotSubdirectory(".snapshots/render_tree"),
		cupaloy.SnapshotFileExtension(".snapshot"),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := &Node{
				NodeType:   "dir",
				RawOutPath: "{{ .Project }}",
			}

			projFS, err := fs.Sub(tt.fs, "{{ .Project }}")
			require.NoError(t, err)

			err = buildNodeTree(projFS, root)
			require.NoError(t, err)

			flat := root.FlattenWriters()

			for _, node := range flat {
				err := RenderNode(tEngine, node, vars)
				require.NoError(t, err)
			}

			snapshot.SnapshotT(t, root.ToAst())
		})
	}
}

var _ RenderableNode = &Node{}

type Node struct {
	NodeType       string
	RawOutPath     string
	WrittenOutPath string
	ReadContent    []byte
	WriteContent   []byte
	Children       []*Node
}

func (n *Node) GetTemplatePath() string { return n.RawOutPath }
func (n *Node) SetOutPath(path string)  { n.WrittenOutPath = path }

func (n *Node) FlattenWriters() []*Node {
	var writers []*Node

	writers = append(writers, n)

	for _, child := range n.Children {
		writers = append(writers, child.FlattenWriters()...)
	}

	return writers
}

func (n *Node) Read(p []byte) (int, error) {
	if n.ReadContent == nil {
		return 0, io.EOF
	}
	retN := copy(p, n.ReadContent)
	n.ReadContent = nil
	return retN, nil
}

func (n *Node) Write(p []byte) (int, error) {
	if n.NodeType == "dir" {
		return len(p), nil
	}

	n.WriteContent = append(n.WriteContent, p...)
	return len(p), nil
}

func (n *Node) childAST(indent int) string {
	ast := &strings.Builder{}

	indentStr := strings.Repeat("\t", indent)

	for _, child := range n.Children {
		println("Child:", child.NodeType, child.RawOutPath, child.WrittenOutPath)
		prefix := fmt.Sprintf("%s:  (type=%s)", child.WrittenOutPath, child.NodeType)
		if child.NodeType == "file" {
			// ensure all new lines in the file are indented
			content := string(bytes.ReplaceAll(child.WriteContent, []byte("\n"), []byte("\n"+indentStr+"\t")))

			ast.WriteString(indentStr + prefix + "\n" + indentStr + "\t" + content + "\n")
		} else {
			ast.WriteString(indentStr + prefix + "\n" + child.childAST(indent+1))
		}
	}

	return ast.String()
}

func (n *Node) ToAst() string {
	rootAST := &strings.Builder{}
	rootAST.WriteString(fmt.Sprintf("%s:  (type=%s)\n", n.WrittenOutPath, n.NodeType))
	rootAST.WriteString(n.childAST(1))
	return rootAST.String()
}

func buildNodeTree(subFs fs.FS, root *Node) error {
	files, err := fs.ReadDir(subFs, ".")
	if err != nil {
		return err
	}

	for _, file := range files {
		node := &Node{
			RawOutPath: file.Name(),
		}
		if file.IsDir() {
			node.NodeType = "dir"
			subsubFS, _ := fs.Sub(subFs, file.Name())

			err = buildNodeTree(subsubFS, node)
			if err != nil {
				return err
			}
		} else {
			node.NodeType = "file"
			readContent, err := fs.ReadFile(subFs, file.Name())
			if err != nil {
				return err
			}
			node.ReadContent = readContent
		}
		root.Children = append(root.Children, node)
	}

	return nil
}
