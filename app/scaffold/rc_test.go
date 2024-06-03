package scaffold

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"strings"
	"testing"

	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewScaffoldRC(t *testing.T) {
	badScaffoldRC := []byte(`---
shorts:
    gh: github.com
    gitea: gitea.com
aliases:
    cli: app/cli
`)

	goodScaffoldRC := []byte(`---
shorts:
    gh: https://github.com
    gitea: https://gitea.com
aliases:
    cli: ~/app/cli
`)

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
			rc, err := NewScaffoldRC(tt.r)
			require.NoError(t, err, "failed on setup")

			err = rc.Validate()

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

func TestScaffoldRC_Authenticator(t *testing.T) {
	authscaffoldRC := []byte(`---
auth:
  - match: "^https://github.com"
    basic:
      username: user
      password: pass
  - match: "^https://gitea.com"
    token: token-value
`)

	rc, err := NewScaffoldRC(bytes.NewReader(authscaffoldRC))
	require.NoError(t, err)

	auth, ok := rc.Authenticator("https://github.com/hay-kot/scaffold-go")
	require.True(t, ok)

	basic, ok := auth.(*githttp.BasicAuth)
	require.True(t, ok)

	assert.Equal(t, "user", basic.Username)
	assert.Equal(t, "pass", basic.Password)

	// Token Auth
	auth, ok = rc.Authenticator("https://gitea.com/hay-kot/scaffold-go")
	require.True(t, ok)

	token, ok := auth.(*githttp.TokenAuth)
	require.True(t, ok)

	assert.Equal(t, "token-value", token.Token)

	// No Auth
	auth, ok = rc.Authenticator("https://gitlab.com/hay-kot/scaffold-go")
	require.False(t, ok)
	assert.Nil(t, auth)
}

func TestScaffoldRC_Authenticator_EnvExpansion(t *testing.T) {
	const (
		ScaffoldGithubUser = "scaffold-gh-user"
		ScaffoldGithubPass = "scaffold-gh-pass"
		ScaffoldGiteaToken = "scaffold-gitea-token"
	)

	t.Setenv("GITHUB_USER", ScaffoldGithubUser)
	t.Setenv("GITHUB_PASS", ScaffoldGithubPass)
	t.Setenv("GITEA_AUTH_TOKEN", ScaffoldGiteaToken)

	authscaffoldRC := []byte(`---
auth:
  - match: "^https://github.com"
    basic:
      username: ${GITHUB_USER}
      password: ${GITHUB_PASS}
  - match: "^https://gitea.com"
    token: ${GITEA_AUTH_TOKEN}
`)

	rc, err := NewScaffoldRC(bytes.NewReader(authscaffoldRC))
	require.NoError(t, err)

	auth, ok := rc.Authenticator("https://github.com/hay-kot/scaffold-go")
	require.True(t, ok)

	basic, ok := auth.(*githttp.BasicAuth)
	require.True(t, ok)

	assert.Equal(t, ScaffoldGithubUser, basic.Username)
	assert.Equal(t, ScaffoldGithubPass, basic.Password)

	// Token Auth
	auth, ok = rc.Authenticator("https://gitea.com/hay-kot/scaffold-go")
	require.True(t, ok)

	token, ok := auth.(*githttp.TokenAuth)
	require.True(t, ok)

	assert.Equal(t, ScaffoldGiteaToken, token.Token)

	// No Auth
	auth, ok = rc.Authenticator("https://gitlab.com/hay-kot/scaffold-go")
	require.False(t, ok)
	assert.Nil(t, auth)
}

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
