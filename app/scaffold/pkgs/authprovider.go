package pkgs

import "github.com/go-git/go-git/v5/plumbing/transport"

type AuthProvider interface {
	Authenticator(pkgurl string) (auth transport.AuthMethod, ok bool)
}

type AuthProviderFunc func(pkgurl string) (auth transport.AuthMethod, ok bool)

func (f AuthProviderFunc) Authenticator(pkgurl string) (auth transport.AuthMethod, ok bool) {
	return f(pkgurl)
}
