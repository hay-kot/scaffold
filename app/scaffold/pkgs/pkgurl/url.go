// Package pkgurl contains functions for parsing remote urls
package pkgurl

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
	"strings"
)

var (
	isSchemeRegExp = regexp.MustCompile(`^[^:]+://`)

	// Ref: https://github.com/git/git/blob/master/Documentation/urls.txt#L37
	scpLikeURLRegExp = regexp.MustCompile(`^(?:(?P<user>[^@]+)@)?(?P<host>[^:\s]+):(?:(?P<port>[0-9]{1,5}):)?(?P<path>[^\\][^#]*)(#(?P<hash>.*))?$`)
)

type Parts struct {
	GitUser string
	Host    string
	Port    string
	Path    string

	// Fragment is the part of the URL after the hash (#)
	Fragment string
	Version  string
}

func (p Parts) PathParts() []string {
	return strings.Split(p.Path, "/")
}

// RepoOwnerAndName returns the owner and name of the repository if the URL format contains exactly two parts. Otherwise,
// it returns false.
//
// Example:
//   - Path="hay-kot/scaffold-go-cli" => user="hay-kot", repo="scaffold-go-cli", ok=true
//   - Path="hay-kot" => user="", repo="", ok=false
func (p Parts) RepoOwnerAndName() (user, repo string, ok bool) {
	parts := p.PathParts()
	if len(parts) < 2 {
		return "", "", false
	}

	return parts[0], parts[1], true
}

func Parse(strURL string) (Parts, error) {
	if MatchesScheme(strURL) {
		// traditional http url
		// Split the url by the hashbang separator, if it exists
		urlParts := strings.Split(strURL, "#")

		url, err := url.ParseRequestURI(urlParts[0])
		if err != nil {
			return Parts{}, fmt.Errorf("failed to parse url: %w", err)
		}

		host := url.Host
		split := strings.Split(url.Path, "/")
		fragment := ""
		if len(urlParts) > 1 {
			fragment = urlParts[1]
		}

		// Remove .git from repo name if it exists but keeps @tag or @branch intact
		split[len(split)-1] = strings.Replace(split[len(split)-1], ".git", "", 1)

		if len(split) < 3 {
			return Parts{}, errors.New("invalid url")
		}

		user := split[1]
		repo := split[2]

		version := ""
		if strings.Contains(repo, "@") {
			split := strings.Split(repo, "@")
			if len(split) != 2 {
				return Parts{}, errors.New("invalid url, unable to parse version")
			}

			repo = split[0]
			version = split[1]
		}

		return Parts{
			GitUser:  user,
			Host:     host,
			Port:     "",
			Path:     user + "/" + repo,
			Fragment: fragment,
			Version:  version,
		}, nil
	}

	if MatchesScpLike(strURL) {
		m := scpLikeURLRegExp.FindStringSubmatch(strURL)

		path := strings.TrimSuffix(m[4], ".git")

		version := ""
		if strings.Contains(path, "@") {
			split := strings.Split(path, "@")
			if len(split) != 2 {
				return Parts{}, errors.New("invalid url, unable to parse version")
			}

			path = split[0]
			version = split[1]
		}

		return Parts{
			GitUser:  m[1],
			Host:     m[2],
			Port:     m[3],
			Path:     path,
			Fragment: m[6],
			Version:  version,
		}, nil
	}

	return Parts{}, errors.New("failed to parse url: matches neither scheme nor scp-like url structure")
}

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
func FindScpLikeComponents(url string) (user, host, port, path, hash string) {
	m := scpLikeURLRegExp.FindStringSubmatch(url)
	return m[1], m[2], m[3], m[4], m[6]
}

// IsLocalEndpoint returns true if the given URL string specifies a
// local file endpoint.  For example, on a Linux machine,
// `/home/user/src/go-git` would match as a local endpoint, but
// `https://github.com/src-d/go-git` would not.
func IsLocalEndpoint(url string) bool {
	return !IsRemoteEndpoint(url)
}
