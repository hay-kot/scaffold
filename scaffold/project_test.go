package scaffold

import (
	"io/fs"
	"testing"

	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadProject(t *testing.T) {
	type args struct {
		fs fs.FS
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "invalid project",
			args: args{
				fs: InvalidProject(),
			},
			wantErr: true,
		},
		{
			name: "valid project",
			args: args{
				fs: DynamicFiles(),
			},
		},
		{
			name: "nested project",
			args: args{
				fs: NestedFiles(),
			},
		},
	}

	snapshot := cupaloy.New(
		cupaloy.SnapshotSubdirectory(".snapshots/project_tree"),
		cupaloy.SnapshotFileExtension(".snapshot"),
	)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadProject(tt.args.fs)

			switch {
			case tt.wantErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				assert.NotNil(t, got)

				assert.NotNil(t, got.Tree)
				snapshot.SnapshotT(t, got.Tree)
			}
		})
	}
}
