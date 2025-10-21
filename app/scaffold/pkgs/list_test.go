package pkgs

import (
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestListFromFS(t *testing.T) {
	tests := []struct {
		name    string
		fs      fstest.MapFS
		want    []string
		wantErr bool
	}{
		{
			name: "empty filesystem",
			fs:   fstest.MapFS{},
			want: []string{},
		},
		{
			name: "single scaffold with yaml",
			fs: fstest.MapFS{
				"cli/scaffold.yaml": &fstest.MapFile{Data: []byte("name: cli")},
			},
			want: []string{"cli"},
		},
		{
			name: "single scaffold with yml",
			fs: fstest.MapFS{
				"api/scaffold.yml": &fstest.MapFile{Data: []byte("name: api")},
			},
			want: []string{"api"},
		},
		{
			name: "multiple scaffolds",
			fs: fstest.MapFS{
				"cli/scaffold.yaml":        &fstest.MapFile{Data: []byte("name: cli")},
				"api/scaffold.yaml":        &fstest.MapFile{Data: []byte("name: api")},
				"microservice/scaffold.yml": &fstest.MapFile{Data: []byte("name: microservice")},
			},
			want: []string{"cli", "api", "microservice"},
		},
		{
			name: "scaffold in nested directory",
			fs: fstest.MapFS{
				"project/backend/scaffold.yaml": &fstest.MapFile{Data: []byte("name: backend")},
			},
			want: []string{"backend"},
		},
		{
			name: "no scaffold files",
			fs: fstest.MapFS{
				"cli/README.md": &fstest.MapFile{Data: []byte("# CLI")},
				"api/docs.txt":  &fstest.MapFile{Data: []byte("docs")},
			},
			want: []string{},
		},
		{
			name: "mixed files and scaffolds",
			fs: fstest.MapFS{
				"cli/scaffold.yaml": &fstest.MapFile{Data: []byte("name: cli")},
				"cli/README.md":     &fstest.MapFile{Data: []byte("# CLI")},
				"api/scaffold.yaml": &fstest.MapFile{Data: []byte("name: api")},
				"docs/README.md":    &fstest.MapFile{Data: []byte("# Docs")},
			},
			want: []string{"cli", "api"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListFromFS(tt.fs)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}
