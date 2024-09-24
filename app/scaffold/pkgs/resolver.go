package pkgs

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/rs/zerolog/log"
)

var ErrNoMatchingScaffold = fmt.Errorf("no matching scaffold")

type ResolverOption func(*Resolver)

func WithCloner(cloner Cloner) ResolverOption {
	return func(r *Resolver) {
		r.cloner = cloner
	}
}

type Resolver struct {
	cloner Cloner
	shorts map[string]string
	cache  string
	cwd    string
}

var commitHashRegex = regexp.MustCompile(`^[0-9a-f]{40}$`)

func isGitCommitHash(s string) bool {
	return commitHashRegex.MatchString(s)
}

func NewResolver(shorts map[string]string, cache, cwd string, opts ...ResolverOption) *Resolver {
	r := &Resolver{
		shorts: shorts,
		cache:  cache,
		cwd:    cwd,
		cloner: ClonerFunc(func(path string, version string, isBare bool, cfg *git.CloneOptions) (string, error) {
			// Clone Repository to cache and set path to cache path
			r, err := git.PlainClone(path, isBare, cfg)
			if err != nil {
				return "", fmt.Errorf("failed to clone repository: %w", err)
			}

			// Get cloned repository path
			wt, err := r.Worktree()
			if err != nil {
				return "", fmt.Errorf("failed to get worktree: %w", err)
			}

			if version != "" && strings.ToLower(version) != "head" {
				if isGitCommitHash(version) {
					err := wt.Checkout(&git.CheckoutOptions{
						Hash:  plumbing.NewHash(version),
						Force: true,
					})
					if err == nil {
						return wt.Filesystem.Root(), nil
					}
				} else {
					// try checkout tag
					tag, err := r.Tag(version)
					if err == nil {
						err := wt.Checkout(&git.CheckoutOptions{
							Hash:  tag.Hash(),
							Force: true,
						})
						if err == nil {
							return wt.Filesystem.Root(), nil
						}
					}

					// try checkout branch
					err = wt.Checkout(&git.CheckoutOptions{
						Branch: plumbing.NewBranchReferenceName(version),
						Force:  true,
					})
					if err != nil {
						return "", fmt.Errorf("failed to checkout branch/tag '%s': %w", version, err)
					}
				}
			}

			return wt.Filesystem.Root(), nil
		}),
	}

	for _, opt := range opts {
		opt(r)
	}

	return r
}

func (r *Resolver) Resolve(arg string, checkDirs []string, authprovider AuthProvider) (path string, err error) {
	remoteRef, isRemote := IsRemote(arg, r.shorts)

	switch {
	case isRemote:
		path, err = r.resolveRemote(remoteRef, authprovider)
		if err != nil {
			return "", fmt.Errorf("failed to resolve remote path: %w", err)
		}
	case filepath.IsAbs(arg):
		path, err = r.resolveAbsolute(arg)
		if err != nil {
			return "", fmt.Errorf("failed to resolve absolute path: %w", err)
		}
	case strings.Contains(arg, "/"):
		path, err = r.resolveRelative(arg)
		if err != nil {
			return "", fmt.Errorf("failed to resolve relative path: %w", err)
		}
	default:
		path, err = r.resolveCwd(arg, checkDirs)
		if err != nil {
			return "", ErrNoMatchingScaffold
		}
	}

	_, err = os.Stat(path)
	if err != nil {
		return "", ErrNoMatchingScaffold
	}

	return path, nil
}

func (r *Resolver) resolveRemote(remoteRef string, authprovider AuthProvider) (path string, err error) {
	pkg, err := ParseRemote(remoteRef)
	if err != nil {
		return "", fmt.Errorf("failed to parse path: %w", err)
	}

	cloneDir := pkg.CloneDir(r.cache)

	_, err = os.Stat(cloneDir)

	remoteRef = strings.TrimSuffix(remoteRef, pkg.Version)
	remoteRef = strings.TrimSuffix(remoteRef, "@")

	switch {
	case err == nil:
		path = pkg.ScaffoldDir(r.cache)
	case os.IsNotExist(err):
		cfg := &git.CloneOptions{
			URL:      remoteRef,
			Progress: os.Stdout,
			Auth:     nil,
			Tags:     git.AllTags,
		}

		auth, ok := authprovider.Authenticator(remoteRef)
		if ok {
			log.Debug().
				Str("url", remoteRef).
				Str("provider", auth.Name()).
				Msg("matching auth provider found")
			cfg.Auth = auth
		}

		// Clone Repository to cache and set path to cache path
		_, err := r.cloner.Clone(cloneDir, pkg.Version, false, cfg)
		if err != nil {
			// ensure directory is cleaned up
			_ = os.RemoveAll(cloneDir)

			return "", fmt.Errorf("failed to clone repository: %w", err)
		}

		path = pkg.ScaffoldDir(r.cache)
	default:
		return "", fmt.Errorf("failed to check if repository is cached: %w", err)
	}

	return path, nil
}

func (r *Resolver) resolveAbsolute(arg string) (string, error) {
	return arg, nil
}

func (r *Resolver) resolveRelative(arg string) (path string, err error) {
	return filepath.Join(r.cwd, arg), nil
}

func (r *Resolver) resolveCwd(arg string, checkDirs []string) (path string, err error) {
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

	if path == "" {
		return "", ErrNoMatchingScaffold
	}

	return path, nil
}
