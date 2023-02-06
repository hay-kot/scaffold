package scaffold

import (
	"io/fs"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/hay-kot/scaffold/internal/core/rwfs"
	"github.com/hay-kot/scaffold/internal/engine"
	"github.com/stretchr/testify/require"
)

func Test_RenderRWFileSystem(t *testing.T) {
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

	vars := engine.Vars{
		"Name":  "Your Name1",
		"Name2": "Your Name2",
	}

	snapshot := cupaloy.New(
		cupaloy.SnapshotSubdirectory(".snapshots/render_rwfs"),
		cupaloy.SnapshotFileExtension(".snapshot"),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			memFS := rwfs.NewMemoryWFS()

			args := &RWFSArgs{
				ReadFS:  tt.fs,
				WriteFS: memFS,
				Project: &Project{
					NameTemplate: "{{ .Project }}",
					Name:         "NewProject",
				},
			}

			root := &AstNode{
				NodeType: DirNodeType,
				Path:     "ROOT_NODE",
			}

			err := RenderRWFS(tEngine, args, vars)
			require.NoError(t, err)

			err = buildNodeTree(memFS, root)
			require.NoError(t, err)

			snapshot.SnapshotT(t, root.String())
		})
	}
}
