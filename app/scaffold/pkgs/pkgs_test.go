package pkgs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParsePath(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr bool
	}{
		{
			name:    "empty",
			input:   "",
			want:    "",
			wantErr: true,
		},
		{
			name:    "no slash",
			input:   "foo",
			want:    "foo",
			wantErr: true,
		},
		{
			name:  "github url",
			input: "https://github.com/hay-kot/scaffold",
			want:  "github.com/hay-kot/scaffold",
		},
		{
			name:  "github url with .git",
			input: "https://github.com/hay-kot/scaffold.git",
			want:  "github.com/hay-kot/scaffold",
		},
		{
			name:  "github url with .git",
			input: "https://github.com/hay-kot/scaffold.git@1.0.2",
			want:  "github.com/hay-kot/scaffold@1.0.2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRemote(tt.input)

			switch {
			case tt.wantErr:
				assert.Error(t, err)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}

func TestIsRemote(t *testing.T) {
	type args struct {
		str    string
		shorts map[string]string
	}
	tests := []struct {
		name         string
		args         args
		wantExpanded string
		wantOk       bool
	}{
		{
			name: "empty",
			args: args{
				str:    "",
				shorts: map[string]string{},
			},
			wantOk: false,
		},
		{
			name: "no slash",
			args: args{
				str:    "foo",
				shorts: map[string]string{},
			},
			wantOk: false,
		},
		{
			name: "github url",
			args: args{
				str:    "gh:hay-kot/scaffold",
				shorts: map[string]string{"gh": "https://github.com"},
			},
			wantExpanded: "https://github.com/hay-kot/scaffold",
			wantOk:       true,
		},
		{
			name: "github url with .git",
			args: args{
				str:    "https://github.com/hay-kot/scaffold.git",
				shorts: map[string]string{},
			},
			wantExpanded: "https://github.com/hay-kot/scaffold.git",
			wantOk:       true,
		},
	}
	for _, tt := range tests {
		expanded, isRemote := IsRemote(tt.args.str, tt.args.shorts)

		switch tt.wantOk {
		case true:
			assert.True(t, isRemote)
			assert.Equal(t, tt.wantExpanded, expanded)
		default:
			assert.False(t, isRemote)
		}
	}
}

func Test_cleanRemoteURL(t *testing.T) {
	type tcase struct {
		name  string
		input string
		want  string
	}

	cases := []tcase{
		{
			name:  "github url (https)",
			input: "https://github.com/hay-kot/scaffold",
			want:  "github.com/hay-kot/scaffold",
		},
		{
			name:  "github url (http)",
			input: "http://github.com/hay-kot/scaffold",
			want:  "github.com/hay-kot/scaffold",
		},
		{
			name:  "github url with .git",
			input: "https://github.com/hay-kot/scaffold.git",
			want:  "github.com/hay-kot/scaffold",
		},
		{
			name:  "github url with ssh prefix",
			input: "git@github.com:hay-kot/scaffold.git",
			want:  "github.com/hay-kot/scaffold",
		},
		{
			name:  "filepath",
			input: "/path/to/repo",
			want:  "/path/to/repo",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := cleanRemoteURL(tc.input)
			assert.Equal(t, tc.want, got)
		})
	}
}
