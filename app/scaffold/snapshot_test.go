package scaffold

import (
	"io/fs"
	"maps"
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
		{
			name: "injected",
			fs:   InjectedFiles(),
			p:    InjectedFilesProject(),
		},
		{
			name: "custom delims",
			fs:   CustomDelimsFiles(),
			p:    CustomDelimsProject(),
		},
		{
			name: "partials",
			fs:   PartialsFiles(),
			p:    PartialsProject(),
		},
		{
			name: "skip with rewrite",
			fs:   SkipWithRewriteFiles(),
			p:    SkipWithRewriteProject(),
		},
		{
			name: "template scaffold",
			fs:   TemplateScaffoldFiles(),
			p:    TemplateScaffoldProject(),
			vars: engine.Vars{
				"name": "MyComponent",
			},
		},
		{
			name: "each expand (directory)",
			fs:   EachExpandFiles(),
			p:    EachExpandProject(),
			vars: engine.Vars{
				"services": []string{"auth", "users"},
			},
		},
		{
			name: "each expand (file)",
			fs:   EachExpandFileFiles(),
			p:    EachExpandFileProject(),
			vars: engine.Vars{
				"items": []string{"foo", "bar", "baz"},
			},
		},
		{
			name: "each expand (as)",
			fs:   EachExpandAsFiles(),
			p:    EachExpandAsProject(),
			vars: engine.Vars{
				"models": []string{"user", "blog_post"},
			},
		},
		{
			name: "each expand (literal)",
			fs:   EachExpandLiteralFiles(),
			p:    EachExpandLiteralProject(),
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
			maps.Copy(vars, tt.vars)
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

			vars, err := BuildVars(tEngine, args.Project, vars)
			require.NoError(t, err)

			err = RenderRWFS(tEngine, args, vars)
			require.NoError(t, err)

			err = fsast.Build(memFS, root)
			require.NoError(t, err)

			snapshot.SnapshotT(t, root.String())
		})
	}
}
