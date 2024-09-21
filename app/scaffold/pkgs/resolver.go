package pkgs

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
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

func NewResolver(shorts map[string]string, cache, cwd string, opts ...ResolverOption) *Resolver {
	r := &Resolver{
		shorts: shorts,
		cache:  cache,
		cwd:    cwd,
		cloner: ClonerFunc(func(path string, isBare bool, cfg *git.CloneOptions) (string, error) {
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
	parsedPath, subdir, err := ParseRemote(remoteRef)
	if err != nil {
		return "", fmt.Errorf("failed to parse path: %w", err)
	}

	dir := filepath.Join(r.cache, parsedPath)

	_, err = os.Stat(dir)

	switch {
	case err == nil:
		path = filepath.Join(dir, subdir)
	case os.IsNotExist(err):
		cfg := &git.CloneOptions{
			URL:      remoteRef,
			Progress: os.Stdout,
			Auth:     nil,
		}

		auth, ok := authprovider.Authenticator(remoteRef)
		if ok {
			log.Debug().Msg("matching auth provider found")
			cfg.Auth = auth
		}

		// Clone Repository to cache and set path to cache path
		clonedPath, err := r.cloner.Clone(dir, false, cfg)
		if err != nil {
			// ensure directory is cleaned up
			_ = os.RemoveAll(dir)

			return "", fmt.Errorf("failed to clone repository: %w", err)
		}

		path = filepath.Join(clonedPath, subdir)
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
