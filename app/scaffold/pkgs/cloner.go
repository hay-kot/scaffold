package pkgs

import "github.com/go-git/go-git/v5"

type Cloner interface {
	Clone(path string, isBare bool, cfg *git.CloneOptions) (string, error)
}

type ClonerFunc func(path string, isBare bool, cfg *git.CloneOptions) (string, error)

func (f ClonerFunc) Clone(path string, isBare bool, cfg *git.CloneOptions) (string, error) {
	return f(path, isBare, cfg)
}
