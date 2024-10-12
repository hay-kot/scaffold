// Package pkgs contains functions for parsing remote urls and checking if a
// directory is a git repository.
package pkgs

import (
	"errors"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs/pkgurl"
)

type Package struct {
	relativePath string

	RemoteRef string
	Version   string
	Subdir    string
}

// CloneDir returns the directory where the repository will be cloned.
func (p Package) CloneDir(root string) string {
	versionstr := ""
	if p.Version != "" && strings.EqualFold(p.Version, "head") {
		versionstr = "@" + p.Version
	}

	return filepath.Join(root, p.relativePath+versionstr)
}

// ScaffoldDir returns the directory where the repository will be cloned and the subdirectory to use if the
// package reference was specified with a subdirectory argument.
func (p Package) ScaffoldDir(root string) string {
	return filepath.Join(p.CloneDir(root), p.Subdir)
}

// ParseRemote parses a remote endpoint and returns a filesystem path representing the
// repository. In addition, it returns the subdirectory of the repository to be used, if any.
//
// Examples:
//
//		ParseRemote("https://github.com/hay-kot/scaffold-go-cli")
//		github.com
//		└── hay-kot
//		    └── scaffold-go-cli
//				     └── repository files
//
//		ParseRemote("https://github.com/hay-kot/scaffold-go-cli#subdir")
//		github.com
//		└── hay-kot
//		    └── scaffold-go-cli
//			 	    └── subdir
//	              └── files to use
func ParseRemote(urlStr string) (Package, error) {
	parts, err := pkgurl.Parse(urlStr)
	if err != nil {
		return Package{}, err
	}

	user, repo, ok := parts.RepoOwnerAndName()
	if !ok {
		return Package{}, errors.New("invalid url, could not identify user and/or repo names")
	}

	relPath := filepath.Join(parts.Host, user, repo)

	pkg := Package{
		RemoteRef:    urlStr,
		relativePath: relPath,
		Version:      parts.Version,
		Subdir:       parts.Fragment,
	}

	return pkg, nil
}

// IsRemote checks if the string is a remote url or an alias for a remote url
// if it is a remote url, it returns the url. If the string uses and alias
// it returns the expanded url
//
// Examples:
//
//	IsRemote(gh:foo/bar) -> https://github.com/foo/bar, true
func IsRemote(str string, shorts map[string]string) (expanded string, ok bool) {
	split := strings.Split(str, ":")

	if len(split) == 2 {
		short := split[0]

		for k, v := range shorts {
			if k == short {
				out, err := url.JoinPath(v, split[1])
				if err != nil {
					return "", false
				}

				return out, pkgurl.IsRemoteEndpoint(out)
			}
		}
	}

	if pkgurl.IsRemoteEndpoint(str) {
		return str, true
	}

	return "", false
}
