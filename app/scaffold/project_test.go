package scaffold

import (
	"io/fs"
	"testing"
	"testing/fstest"

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
				fs: InvalidStructure(),
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
			got, err := LoadProject(tt.args.fs, Options{})

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

func TestRootDir(t *testing.T) {
	file := func() *fstest.MapFile { return &fstest.MapFile{Data: []byte("x")} }

	tests := []struct {
		name    string
		fs      fstest.MapFS
		want    string
		wantErr string
	}{
		{
			name: "project dir",
			fs:   fstest.MapFS{"{{ .ProjectKebab }}/a.txt": file()},
			want: "{{ .ProjectKebab }}",
		},
		{
			name: "templates dir",
			fs:   fstest.MapFS{"templates/a.txt": file()},
			want: TemplateDirName,
		},
		{
			name:    "no root dir",
			fs:      fstest.MapFS{"other/a.txt": file()},
			wantErr: "does not exist",
		},
		{
			name:    "templates file is not a root dir",
			fs:      fstest.MapFS{"templates": file()},
			wantErr: "does not exist",
		},
		{
			name: "project and templates dirs",
			fs: fstest.MapFS{
				"{{ .Project }}/a.txt": file(),
				"templates/b.txt":      file(),
			},
			wantErr: `found 2: "{{ .Project }}", "templates"`,
		},
		{
			name: "two project dirs",
			fs: fstest.MapFS{
				"{{ .Project }}/a.txt":      file(),
				"{{ .ProjectSnake }}/b.txt": file(),
			},
			wantErr: `found 2: "{{ .Project }}", "{{ .ProjectSnake }}"`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := RootDir(tt.fs)
			if tt.wantErr != "" {
				require.ErrorContains(t, err, tt.wantErr)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestLoadProject_RejectsMultipleRootDirs(t *testing.T) {
	fsys := fstest.MapFS{
		"scaffold.yaml":        &fstest.MapFile{Data: []byte("questions: []\n")},
		"{{ .Project }}/a.txt": &fstest.MapFile{Data: []byte("a")},
		"templates/b.txt":      &fstest.MapFile{Data: []byte("b")},
	}

	_, err := LoadProject(fsys, Options{})
	require.ErrorContains(t, err, "exactly one root directory")
}
