package scaffold

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInject(t *testing.T) {
	const Marker = "# Inject Marker"
	const Input = `---
hello world
    indented line
    # Inject Marker
`

	type args struct {
		s    string
		data string
		at   string
		mode Mode
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "inject",
			args: args{
				s:    Input,
				data: "injected line 1\ninjected line 2",
				at:   Marker,
			},
			want: `---
hello world
    indented line
    injected line 1
    injected line 2
    # Inject Marker
`,
		},
		{
			name: "inject after",
			args: args{
				s:    Input,
				data: "injected line 1\ninjected line 2",
				at:   Marker,
				mode: After,
			},
			want: `---
hello world
    indented line
    # Inject Marker
    injected line 1
    injected line 2
`,
		},
		{
			name: "inject no marker",
			args: args{
				s:    Input,
				data: "injected",
				at:   Marker + "invalid",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Inject(strings.NewReader(tt.args.s), tt.args.data, tt.args.at, tt.args.mode)

			switch {
			case tt.wantErr:
				assert.Error(t, err)
			default:
				assert.NoError(t, err)
				assert.Equal(t, tt.want, string(got))
			}
		})
	}
}
