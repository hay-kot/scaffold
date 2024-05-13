package scaffold

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var t1 = `---
hello world
    indented line
    # Inject Marker
`

var t1Want = `---
hello world
    indented line
    injected line 1
    injected line 2
    # Inject Marker
`

var t2 = `---
hello world
    indented line
    # Inject After Marker
`

var t2Want = `---
hello world
    indented line
    # Inject After Marker
    injected line 1
    injected line 2
`

func TestInject(t *testing.T) {
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
				s:    t1,
				data: "injected line 1\ninjected line 2",
				at:   "# Inject Marker",
			},
			want: t1Want,
		},
		{
			name: "inject after",
			args: args{
				s:    t2,
				data: "injected line 1\ninjected line 2",
				at:   "# Inject After Marker",
				mode: After,
			},
			want: t2Want,
		},
		{
			name: "inject no marker",
			args: args{
				s:    t1,
				data: "injected line 1\ninjected line 2",
				at:   "# Inject Marker 2",
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
