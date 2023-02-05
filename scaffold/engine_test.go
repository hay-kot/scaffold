package scaffold

import (
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	tEngine = NewEngine()
)

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
			tmpl: "./path/to/file/{{ .Scaffold.Name }}",
			want: "./path/to/file/Name",
			vars: Vars{
				"Name": "Name",
			},
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
			tmpl, err := tEngine.TmplFactory(tt.reader)

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
			tmpl:    "Hello {{ .Scaffold.World }}",
			want:    "Hello World!",
			wantErr: false,
			vars: Vars{
				"World": "World!",
			},
		},
		{
			name:    "Basic template with sprig function",
			tmpl:    "Hello {{ .Scaffold.World | upper }}",
			want:    "Hello WORLD!",
			wantErr: false,
			vars: Vars{
				"World": "world!",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpl, err := tEngine.TmplFactory(strings.NewReader(tt.tmpl))
			require.NoError(t, err, "failed to create template during render test setup")

			strBuff := &strings.Builder{}
			err = tEngine.RenderTemplate(strBuff, tmpl, tt.vars)

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
