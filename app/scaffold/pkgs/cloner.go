package pkgs

import (
	"context"

	"github.com/go-git/go-git/v5"
)

type Cloner interface {
	Clone(ctx context.Context, path string, version string, isBare bool, cfg *git.CloneOptions) (string, error)
}

type ClonerFunc func(ctx context.Context, path string, version string, isBare bool, cfg *git.CloneOptions) (string, error)

func (f ClonerFunc) Clone(ctx context.Context, path string, version string, isBare bool, cfg *git.CloneOptions) (string, error) {
	return f(ctx, path, version, isBare, cfg)
}
