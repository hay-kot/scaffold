// Package scaffoldrc contains the runtime configuration for users
package scaffoldrc

import (
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5/plumbing/transport"
	githttp "github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/hay-kot/scaffold/internal/styles"
	"gopkg.in/yaml.v3"
)

type ScaffoldRC struct {
	// Settings define the settings for the scaffold application.
	Settings Settings `yaml:"settings"`

	// Defaults define a default value for a variable.
	//   name: myproject
	//   git_user: hay-kot
	//
	// These are injected into the template as variables for
	// every scaffold.
	Defaults map[string]any `yaml:"defaults"`

	// Aliases define a alias for a repository.
	// or filepath.
	//
	//   component: https://githublcom/hay-kot/scaffold-go-component
	//   cli: https://github.com/hay-kot/scaffold-go-cli
	Aliases map[string]string `yaml:"aliases"`

	// Shorts define a short name for a repository.
	// the key will be expanded into the value.
	//   gh: https://github.com
	//   gt: https://gitea.com
	//
	// This will allow you to use the short name in the scaffold
	//   gh:myorg/myrepo
	Shorts map[string]string `yaml:"shorts"`

	// Auth defines a list of auth entries that can be used to
	// authenticate with a remote SCM.
	Auth []AuthEntry `yaml:"auth"`
}

type Settings struct {
	Theme    styles.HuhTheme `yaml:"theme"`
	RunHooks RunHooksOption  `yaml:"run_hooks"`
}

type AuthEntry struct {
	Match regexp.Regexp `yaml:"match"`
	Basic BasicAuth     `yaml:"basic"`
	Token string        `yaml:"token"`
}

type BasicAuth struct {
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

// Default returns a default scaffold rc file.
func Default() *ScaffoldRC {
	return &ScaffoldRC{
		Settings: Settings{
			Theme:    styles.HuhThemeScaffold,
			RunHooks: RunHooksPrompt,
		},
	}
}

// New reads a scaffold rc file from the reader and returns a
// ScaffoldRC struct.
func New(r io.Reader) (*ScaffoldRC, error) {
	out := Default()
	err := yaml.NewDecoder(r).Decode(&out)
	if err != nil {
		if errors.Is(err, io.EOF) {
			// Assume empty file and return empty struct
			return out, nil
		}
		return nil, err
	}

	return out, nil
}

func (rc *ScaffoldRC) Validate() error {
	errs := []RCValidationError{}

	for k, v := range rc.Shorts {
		_, err := url.ParseRequestURI(v)
		if err != nil {
			errs = append(errs, RCValidationError{
				Key:   k,
				Cause: fmt.Errorf("parse url failed: %w", err),
			})
		}
	}

	for k, v := range rc.Aliases {
		// Shorts must be absolute path or relative to ~ or a URL
		_, err := url.ParseRequestURI(v)
		if err != nil {
			if !filepath.IsAbs(v) && !strings.HasPrefix(v, "~") {
				errs = append(errs, RCValidationError{
					Key:   k,
					Cause: fmt.Errorf("invalid short path: %w", err),
				})
			}
		}
	}

	if !rc.Settings.Theme.IsValid() {
		errs = append(errs, RCValidationError{
			Key:   "settings.theme",
			Cause: fmt.Errorf("invalid theme: %s", rc.Settings.Theme.String()),
		})
	}

	if !rc.Settings.RunHooks.IsValid() {
		errs = append(errs, RCValidationError{
			Key:   "settings.run_hooks",
			Cause: fmt.Errorf("invalid run_hooks: %s", rc.Settings.RunHooks.String()),
		})
	}

	if len(errs) > 0 {
		return RcValidationErrors(errs)
	}

	return nil
}

func (rc *ScaffoldRC) Authenticator(pkgurl string) (transport.AuthMethod, bool) {
	for _, auth := range rc.Auth {
		if auth.Match.MatchString(pkgurl) {
			if auth.Basic.Username != "" {
				return &githttp.BasicAuth{
					Username: expandEnvVars(auth.Basic.Username),
					Password: expandEnvVars(auth.Basic.Password),
				}, true
			}

			if auth.Token != "" {
				return &githttp.TokenAuth{
					Token: expandEnvVars(auth.Token),
				}, true
			}
		}
	}

	return nil, false
}

// expandEnvVars will expand an environment variable in the form of ${VAR}
func expandEnvVars(s string) string {
	if !strings.HasPrefix(s, "${") && !strings.HasSuffix(s, "}") {
		return s
	}

	return os.Getenv(s[2 : len(s)-1])
}
