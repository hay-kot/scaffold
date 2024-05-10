package scaffold

import (
	"io/fs"
	"testing"

	"github.com/bradleyjkemp/cupaloy/v2"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/fsast"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/stretchr/testify/require"
)

func Test_RenderRWFileSystem(t *testing.T) {
	tests := []struct {
		name string
		fs   fs.FS
		p    *Project
		vars engine.Vars
	}{
		{
			name: "basic",
			fs:   DynamicFiles(),
			p:    DynamicFilesProject(),
		},
		{
			name: "nested",
			fs:   NestedFiles(),
			p:    NestedFilesProject(),
		},
		{
			name: "feature flag (off)",
			fs:   FeatureFlagFiles(),
			p:    FeatureFlagProject(),
		},
		{
			name: "feature flag (on)",
			fs:   FeatureFlagFiles(),
			p:    FeatureFlagProject(),
			vars: engine.Vars{
				"feature": "true",
			},
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
		if tt.vars != nil {
			for k, v := range tt.vars {
				vars[k] = v
			}
		}

		t.Run(tt.name, func(t *testing.T) {
			memFS := rwfs.NewMemoryWFS()

			args := &RWFSArgs{
				ReadFS:  tt.fs,
				WriteFS: memFS,
				Project: tt.p,
			}

			root := &fsast.AstNode{
				NodeType: fsast.DirNodeType,
				Path:     "ROOT_NODE",
			}

			vars, err := BuildVars(tEngine, args, vars)
			require.NoError(t, err)

			err = RenderRWFS(tEngine, args, vars)
			require.NoError(t, err)

			err = fsast.Build(memFS, root)
			require.NoError(t, err)

			snapshot.SnapshotT(t, root.String())
		})
	}
}
