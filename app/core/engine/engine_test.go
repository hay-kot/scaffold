package engine

import (
	"io"
	"strings"
	"testing"
	"testing/fstest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var tEngine = New()

func TestScaffold_TmplString(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		want    string
		wantErr bool
		vars    any
	}{
		{
			name: "Basic template",
			tmpl: "./path/to/file/{{ .Name }}",
			want: "./path/to/file/Name",
			vars: Vars{
				"Name": "Name",
			},
		},
		{
			name: "Test custom func 'wraptmpl'",
			tmpl: "./path/to/file/{{ wraptmpl `Arg` }}",
			want: "./path/to/file/{{ Arg }}",
			vars: Vars{},
		},
		{
			name:    "Empty template",
			tmpl:    "./my/path/without/template",
			want:    "./my/path/without/template",
			wantErr: false,
			vars: Vars{
				"World": "world!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tEngine.TmplString(tt.tmpl, tt.vars)

			switch {
			case tt.wantErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestScaffold_TmplFactory(t *testing.T) {
	tests := []struct {
		name    string
		reader  io.Reader
		wantErr bool
	}{
		{
			name:    "Nil reader",
			reader:  nil,
			wantErr: true,
		},
		{
			name:    "Basic template",
			reader:  strings.NewReader("{{ .Scaffold }}"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := tEngine.Factory(tt.reader)

			switch {
			case tt.wantErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				assert.NotNil(t, tmpl)
			}
		})
	}
}

func TestEngine_RegisterPartialsFS(t *testing.T) {
	partialsFS := fstest.MapFS{
		"header.tmpl":              {Data: []byte("Header")},
		"licence/links.txt":        {Data: []byte("Links for {{ .Name }}")},
		"common/nested/deep.tmpl":  {Data: []byte("Deep")},
		"my-partial.with.dots.txt": {Data: []byte("Dots")},
		".gitkeep":                 {Data: []byte("")},
	}

	e := New()
	require.NoError(t, e.RegisterPartialsFS(partialsFS, "."))

	for _, name := range []string{"header", "licence/links", "common/nested/deep", "my-partial.with.dots", ".gitkeep"} {
		assert.Contains(t, e.partials, name)
	}

	out, err := e.TmplString(`{{ partial "licence/links" . }}`, Vars{"Name": "scaffold"})
	require.NoError(t, err)
	assert.Equal(t, "Links for scaffold", out)
}

func TestEngine_RegisterPartialsFS_Subdirectory(t *testing.T) {
	partialsFS := fstest.MapFS{
		"templates/common/snippet.tmpl": {Data: []byte("Snippet")},
		"outside.tmpl":                  {Data: []byte("Outside")},
	}

	e := New()
	require.NoError(t, e.RegisterPartialsFS(partialsFS, "templates"))

	assert.Contains(t, e.partials, "common/snippet")
	assert.NotContains(t, e.partials, "outside")
}

func TestIsValidPartialName(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{name: "header", want: true},
		{name: "licence/links", want: true},
		{name: "a/b/c", want: true},
		{name: "my-partial.v2", want: true},
		{name: ".gitkeep", want: true},
		{name: "", want: false},
		{name: "/leading", want: false},
		{name: "trailing/", want: false},
		{name: "double//slash", want: false},
		{name: "../escape", want: false},
		{name: "./here", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, IsValidPartialName(tt.name))
		})
	}
}

func TestScaffold_RenderTemplate(t *testing.T) {
	tests := []struct {
		name    string
		tmpl    string
		want    string
		vars    any
		wantErr bool
	}{
		{
			name:    "Basic template",
			tmpl:    "Hello {{ .World }}",
			want:    "Hello World!",
			wantErr: false,
			vars: Vars{
				"World": "World!",
			},
		},
		{
			name:    "Basic template with sprout function",
			tmpl:    "Hello {{ .World | upper }}",
			want:    "Hello WORLD!",
			wantErr: false,
			vars: Vars{
				"World": "world!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := tEngine.Factory(strings.NewReader(tt.tmpl))
			require.NoError(t, err, "failed to create template during render test setup")

			strBuff := &strings.Builder{}
			err = tEngine.Render(strBuff, tmpl, tt.vars)

			switch {
			case tt.wantErr:
				require.Error(t, err)
			default:
				require.NoError(t, err)
				assert.Equal(t, tt.want, strBuff.String())
			}
		})
	}
}
