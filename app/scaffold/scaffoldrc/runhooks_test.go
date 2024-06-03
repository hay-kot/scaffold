package scaffoldrc

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_RunHooksOption_UnmarshalText(t *testing.T) {
	type tcase struct {
		namefmt   string
		inputs    []string
		want      RunHooksOption
		wantValid bool
	}

	tests := []tcase{
		{
			namefmt:   "valid for 'never' (%s)",
			inputs:    []string{"never", "no", "false"},
			want:      RunHooksNever,
			wantValid: true,
		},
		{
			namefmt:   "valid for 'always' (%s)",
			inputs:    []string{"always", "yes", "true"},
			want:      RunHooksAlways,
			wantValid: true,
		},
		{
			namefmt:   "valid for 'prompt' (%s)",
			inputs:    []string{"prompt", ""},
			want:      RunHooksPrompt,
			wantValid: true,
		},
		{
			namefmt:   "invalid for any option (%s)",
			inputs:    []string{"invalid", " "},
			wantValid: false,
		},
	}

	for _, tt := range tests {
		name := fmt.Sprintf(tt.namefmt, strings.Join(tt.inputs, ", "))
		t.Run(name, func(t *testing.T) {
			for _, input := range tt.inputs {
				var got RunHooksOption
				err := got.UnmarshalText([]byte(input))

				require.NoError(t, err) // UnmarshalText should _never_ error

				switch {
				case tt.wantValid:
					assert.Equal(t, tt.want, got)
					assert.True(t, got.IsValid())
				default:
					assert.NotEqual(t, tt.want, got)
					assert.False(t, got.IsValid())
				}
			}
		})
	}
}
