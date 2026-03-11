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
				"cli/scaffold.yaml":         &fstest.MapFile{Data: []byte("name: cli")},
				"api/scaffold.yaml":         &fstest.MapFile{Data: []byte("name: api")},
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

func TestListSystem_SkipsTemplateDirs(t *testing.T) {
	tests := []struct {
		name            string
		fs              fstest.MapFS
		wantRoot        string
		wantSubPackages []string
	}{
		{
			name: "scaffold.yaml inside template dir is not a sub-package",
			fs: fstest.MapFS{
				"github.com/user/repo/.git/HEAD":                                          &fstest.MapFile{Data: []byte("ref: refs/heads/main")},
				"github.com/user/repo/scaffold.yaml":                                      &fstest.MapFile{Data: []byte("name: root")},
				"github.com/user/repo/{{ .Project }}/.scaffold/command/scaffold.yaml":     &fstest.MapFile{Data: []byte("name: command")},
				"github.com/user/repo/{{ .Project }}/.scaffold/command/templates/main.go": &fstest.MapFile{Data: []byte("package main")},
			},
			wantRoot:        "github.com/user/repo",
			wantSubPackages: nil,
		},
		{
			name: "scaffold.yaml inside templates dir is not a sub-package",
			fs: fstest.MapFS{
				"github.com/user/repo/.git/HEAD":                      &fstest.MapFile{Data: []byte("ref: refs/heads/main")},
				"github.com/user/repo/scaffold.yaml":                  &fstest.MapFile{Data: []byte("name: root")},
				"github.com/user/repo/templates/nested/scaffold.yaml": &fstest.MapFile{Data: []byte("name: nested")},
			},
			wantRoot:        "github.com/user/repo",
			wantSubPackages: nil,
		},
		{
			name: "real sub-package is still detected",
			fs: fstest.MapFS{
				"github.com/user/repo/.git/HEAD":               &fstest.MapFile{Data: []byte("ref: refs/heads/main")},
				"github.com/user/repo/scaffold.yaml":           &fstest.MapFile{Data: []byte("name: root")},
				"github.com/user/repo/component/scaffold.yaml": &fstest.MapFile{Data: []byte("name: component")},
			},
			wantRoot:        "github.com/user/repo",
			wantSubPackages: []string{"component"},
		},
		{
			name: "mix of real sub-packages and template dirs",
			fs: fstest.MapFS{
				"github.com/user/repo/.git/HEAD":                                      &fstest.MapFile{Data: []byte("ref: refs/heads/main")},
				"github.com/user/repo/scaffold.yaml":                                  &fstest.MapFile{Data: []byte("name: root")},
				"github.com/user/repo/component/scaffold.yaml":                        &fstest.MapFile{Data: []byte("name: component")},
				"github.com/user/repo/{{ .Project }}/.scaffold/command/scaffold.yaml": &fstest.MapFile{Data: []byte("name: command")},
			},
			wantRoot:        "github.com/user/repo",
			wantSubPackages: []string{"component"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ListSystem(tt.fs)
			require.NoError(t, err)
			require.Len(t, got, 1)
			assert.Equal(t, tt.wantRoot, got[0].Root)
			assert.Equal(t, tt.wantSubPackages, got[0].SubPackages)
		})
	}
}
