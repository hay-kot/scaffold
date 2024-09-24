package pkgs

import "github.com/go-git/go-git/v5"

type Cloner interface {
	Clone(path string, version string, isBare bool, cfg *git.CloneOptions) (string, error)
}

type ClonerFunc func(path string, version string, isBare bool, cfg *git.CloneOptions) (string, error)

func (f ClonerFunc) Clone(path string, version string, isBare bool, cfg *git.CloneOptions) (string, error) {
	return f(path, version, isBare, cfg)
}
