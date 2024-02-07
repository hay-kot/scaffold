package pkgs

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/AlecAivazis/survey/v2"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

var ErrNoMatchingScaffold = fmt.Errorf("no matching scaffold")

type Resolver struct {
	shorts map[string]string
	cache  string
	cwd    string
}

func NewResolver(shorts map[string]string, cache, cwd string) *Resolver {
	return &Resolver{
		shorts: shorts,
		cache:  cache,
		cwd:    cwd,
	}
}

func (r *Resolver) Resolve(arg string, checkDirs []string) (path string, err error) {
	remoteRef, isRemote := IsRemote(arg, r.shorts)

	switch {
	case isRemote:
		parsedPath, err := ParseRemote(remoteRef)
		if err != nil {
			return "", fmt.Errorf("failed to parse path: %w", err)
		}

		dir := filepath.Join(r.cache, parsedPath)

		_, err = os.Stat(dir)

		switch {
		case err == nil:
			path = dir
		case os.IsNotExist(err):
			// Close Repository to cache and set path to cache path
			r, err := git.PlainClone(dir, false, &git.CloneOptions{
				URL:      remoteRef,
				Progress: os.Stdout,
			})
			if err != nil {
				if errors.Is(err, transport.ErrAuthenticationRequired) {
					username, password, err := promptHTTPAuth()
					if err == nil {
						r, err = git.PlainClone(dir, false, &git.CloneOptions{
							URL:      remoteRef,
							Progress: os.Stdout,
							Auth: &http.BasicAuth{
								Username: username,
								Password: password,
							},
						})
					}
					if err != nil {
						return "", fmt.Errorf("failed to clone repository: %w", err)
					}
				} else {
					return "", fmt.Errorf("failed to clone repository: %w", err)
				}
			}

			// Get cloned repository path
			wt, err := r.Worktree()
			if err != nil {
				return "", fmt.Errorf("failed to get worktree: %w", err)
			}

			path = wt.Filesystem.Root()
		default:
			return "", fmt.Errorf("failed to check if repository is cached: %w", err)
		}

	case filepath.IsAbs(arg):
		path = arg
	case strings.Contains(arg, "/"):
		path = filepath.Join(r.cwd, arg)
	default:
		// Otherwise check local .scaffold directory for matching path
		for _, dir := range checkDirs {
			absPath, err := filepath.Abs(filepath.Join(dir, arg))
			if err != nil {
				return "", fmt.Errorf("failed to get absolute path: %w", err)
			}

			// Check if path exists
			_, err = os.Stat(absPath)
			if err == nil {
				path = absPath
				break
			}
		}
	}

	_, err = os.Stat(path)
	if err != nil {
		return "", ErrNoMatchingScaffold
	}

	return path, nil
}

var qs = []*survey.Question{
	{
		Name:     "username",
		Prompt:   &survey.Input{Message: "Username:"},
		Validate: survey.Required,
	},
	{
		Name: "password",
		Prompt: &survey.Password{
			Message: "Password/personal access token:",
		},
	},
}

func promptHTTPAuth() (string, string, error) {
	answers := struct {
		Username string
		Password string
	}{}

	err := survey.Ask(qs, &answers)
	if err != nil {
		return "", "", fmt.Errorf("failed to parse http auth input: %w", err)
	}

	return answers.Username, answers.Password, nil
}
