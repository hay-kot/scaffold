package pkgs

import (
	"fmt"
	"strings"

	"github.com/go-git/go-git/v5"
)

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
		return Version{}, fmt.Errorf("failed to open git repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return Version{}, fmt.Errorf("failed to get HEAD: %w", err)
	}

	repository := path

	// If the repository is a git repository, get the remote url
	remote, err := repo.Remote("origin")
	if err == nil {
		urls := remote.Config().URLs
		if len(urls) > 0 {
			repository = cleanRemoteURL(urls[0])
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
