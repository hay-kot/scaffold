package scaffold

import (
	"io/fs"
	"testing"

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

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadProject(tt.args.fs)

			switch {
			case tt.wantErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				assert.NotNil(t, got)
			}
		})
	}
}
