// Package url contains functions for parsing remote urls
package url

import (
	"regexp"
)

var (
	isSchemeRegExp = regexp.MustCompile(`^[^:]+://`)

	// Ref: https://github.com/git/git/blob/master/Documentation/urls.txt#L37
	scpLikeURLRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5}):)?(?P<path>[^\\].*)$`)
)

// MatchesScheme returns true if the given string matches a URL-like
// format scheme.
func MatchesScheme(url string) bool {
	return isSchemeRegExp.MatchString(url)
}

// MatchesScpLike returns true if the given string matches an SCP-like
// format scheme.
func MatchesScpLike(url string) bool {
	return scpLikeURLRegExp.MatchString(url)
}

// IsRemoteEndpoint returns true if the giver URL string specifies
// a remote endpoint. For example, on a Linux machine,
// `https://github.com/src-d/go-git` would match as a remote
// endpoint, but `/home/user/src/go-git` would not.
func IsRemoteEndpoint(url string) bool {
	return MatchesScheme(url) || MatchesScpLike(url)
}

// FindScpLikeComponents returns the user, host, port and path of the
// given SCP-like URL.
func FindScpLikeComponents(url string) (user, host, port, path string) {
	m := scpLikeURLRegExp.FindStringSubmatch(url)
	return m[1], m[2], m[3], m[4]
}

// IsLocalEndpoint returns true if the given URL string specifies a
// local file endpoint.  For example, on a Linux machine,
// `/home/user/src/go-git` would match as a local endpoint, but
// `https://github.com/src-d/go-git` would not.
func IsLocalEndpoint(url string) bool {
	return !IsRemoteEndpoint(url)
}
