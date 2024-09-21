// Package pkgs contains functions for parsing remote urls and checking if a
// directory is a git repository.
package pkgs

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/hay-kot/scaffold/app/scaffold/pkgs/pkgurl"

	"github.com/go-git/go-git/v5"
)

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
func ParseRemote(urlStr string) (string, string, error) {
	var (
		host   string
		user   string
		repo   string
		subdir string
		err    error
	)

	switch {
	case pkgurl.MatchesScheme(urlStr):
		host, user, repo, subdir, err = parseRemoteURL(urlStr)
	case pkgurl.MatchesScpLike(urlStr):
		host, user, repo, subdir, err = parseRemoteScpLike(urlStr)
	default:
		return "", "", fmt.Errorf("failed to parse url: matches neither scheme nor scp-like url structure")
	}

	if err != nil {
		return "", "", err
	}

	return filepath.Join(host, user, repo), subdir, nil
}

// Parses a remote URL endpoint into its host, user, and repo name
// parts
func parseRemoteURL(urlStr string) (string, string, string, string, error) {
	// Split the url by the hashbang separator, if it exists
	urlParts := strings.Split(urlStr, "#")

	url, err := url.ParseRequestURI(urlParts[0])
	if err != nil {
		return "", "", "", "", fmt.Errorf("failed to parse url: %w", err)
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
		return "", "", "", "", fmt.Errorf("invalid url")
	}

	user := split[1]
	repo := split[2]

	return host, user, repo, fragment, nil
}

// Parses a remote SCP-like endpoint into its host, user, and repo name
// parts
func parseRemoteScpLike(urlStr string) (string, string, string, string, error) {
	user, host, _, path, hash := pkgurl.FindScpLikeComponents(urlStr)

	return host, user, strings.TrimSuffix(path, ".git"), hash, nil
}

// IsRemote checks if the string is a remote url or an alias for a remote url
// if it is a remote url, it returns the url. If the string uses and alias
// it returns the expanded url
//
// Examples:
//
//	isRemote(gh:foo/bar) -> https://github.com/foo/bar, true
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

// Update updates a git repository to the latest commit
func Update(cacheRoot string, scaffoldPath string) (updated bool, err error) {
	// The package may be in a nested directory of the repository, so we need to find the root
	// of the repository

	scaffoldPathParts := strings.Split(scaffoldPath, string(filepath.Separator))

	repoPath := ""
	gitDirFound := false

	// Assuming a minimum of 2 parts in the path i.e. git.host.com/repo
	for i := len(scaffoldPathParts); i > 1; i-- {
		repoPath = filepath.Join(cacheRoot, filepath.Join(scaffoldPathParts[:i]...))
		if _, err := os.Stat(filepath.Join(repoPath, ".git")); err == nil {
			gitDirFound = true
			break
		}
	}

	if !gitDirFound {
		return false, fmt.Errorf("no git repository found for %s", scaffoldPath)
	}

	repo, err := git.PlainOpen(repoPath)
	if err != nil {
		return false, err
	}

	w, err := repo.Worktree()
	if err != nil {
		return false, err
	}

	err = w.Pull(&git.PullOptions{})
	if err != nil {
		if errors.Is(err, git.NoErrAlreadyUpToDate) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

type Version struct {
	Repository string
	Commit     string
	// TODO: add a field for tags once supported
}

func (v Version) String() string {
	return fmt.Sprintf("%s@%s", v.Repository, v.CommitShort())
}

func (v Version) CommitShort() string {
	if len(v.Commit) < 7 {
		return v.Commit
	}

	return v.Commit[:7]
}

func (v Version) IsZero() bool {
	return v.Repository == "" && v.Commit == ""
}

func GetVersion(path string) (Version, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return Version{}, err
	}

	head, err := repo.Head()
	if err != nil {
		return Version{}, err
	}

	repository := path

	// If the repository is a git repository, get the remote url
	remotes, err := repo.Remotes()
	if err == nil && len(remotes) > 0 {
		cfg := remotes[0].Config()
		if len(cfg.URLs) > 0 {
			repository = cleanRemoteURL(cfg.URLs[0])
		}
	}

	return Version{
		Repository: repository,
		Commit:     head.Hash().String(),
	}, nil
}

func cleanRemoteURL(v string) string {
	v = strings.TrimSuffix(v, ".git")
	v = strings.TrimPrefix(v, "https://")
	v = strings.TrimPrefix(v, "http://")
	v = strings.TrimPrefix(v, "git@")

	// handle github.com:hay-kot/scaffold.git
	v = strings.Replace(v, ":", "/", 1)
	return v
}
