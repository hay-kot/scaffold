package commands

import (
	"errors"
	"fmt"

	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/hay-kot/scaffold/app/scaffold/pkgs"
)

func (ctrl *Controller) resolve(
	argPath string,
	outputdir string,
	noPrompt bool,
	force bool,
) (string, error) {
	if argPath == "" {
		return "", fmt.Errorf("path is required")
	}

	// Status() call for go-git is too slow to be used here
	// https://github.com/go-git/go-git/issues/181
	if !force {
		ok := checkWorkingTree(outputdir)
		if !ok {
			return "", fmt.Errorf("working tree is dirty, use --force to override")
		}
	}

	resolver := pkgs.NewResolver(ctrl.rc.Shorts, ctrl.Flags.Cache, ".")

	if v, ok := ctrl.rc.Aliases[argPath]; ok {
		argPath = v
	}

	path, err := resolver.Resolve(argPath, ctrl.Flags.ScaffoldDirs, ctrl.rc)
	if err != nil {
		orgErr := err

		switch {
		case errors.Is(err, transport.ErrAuthenticationRequired):
			if noPrompt {
				return "", err
			}

			username, password, err := httpAuthPrompt(argPath, ctrl.rc.Settings.Theme)
			if err != nil {
				return "", err
			}

			path, err = resolver.Resolve(argPath, ctrl.Flags.ScaffoldDirs, basicAuthAuthorizer(argPath, username, password))
			if err != nil {
				return "", err
			}
		default:
			if noPrompt {
				return "", err
			}

			systemMatches, localMatches, err := ctrl.fuzzyFallBack(argPath)
			if err != nil {
				return "", err
			}

			var first string
			var isSystemMatch bool
			if len(systemMatches) > 0 {
				first = systemMatches[0]
				isSystemMatch = true
			}

			if len(localMatches) > 0 {
				first = localMatches[0]
				isSystemMatch = false
			}

			if first != "" {
				useMatch := didYouMeanPrompt(argPath, first, isSystemMatch)

				if useMatch {
					if isSystemMatch {
						// prepend https:// so it resolves to the correct path
						first = "https://" + first
					}

					resolved, err := resolver.Resolve(first, ctrl.Flags.ScaffoldDirs, ctrl.rc)
					if err != nil {
						return "", err
					}

					path = resolved
				}
			}
		}

		if path == "" {
			return "", fmt.Errorf("failed to resolve path: %w", orgErr)
		}
	}

	return path, nil
}
