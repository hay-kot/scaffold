package scaffold

import (
	"bytes"
	"errors"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var badScaffoldRC = []byte(`---
shorts:
    gh: github.com
    gitea: gitea.com
aliases:
    cli: app/cli
`)

var goodScaffoldRC = []byte(`---
shorts:
    gh: https://github.com
    gitea: https://gitea.com
aliases:
    cli: ~/app/cli
`)

func TestNewScaffoldRC(t *testing.T) {
	tests := []struct {
		name    string
		r       io.Reader
		wantErr bool
	}{
		{
			name:    "bad scaffold rc",
			r:       bytes.NewReader(badScaffoldRC),
			wantErr: true,
		},
		{
			name:    "good scaffold rc",
			r:       bytes.NewReader(goodScaffoldRC),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewScaffoldRC(tt.r)

			switch {
			case tt.wantErr:
				require.Error(t, err)
				assert.True(t, errors.As(err, &RcValidationErrors{}))
			default:
				assert.NoError(t, err)
			}
		})
	}
}
