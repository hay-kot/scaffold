package pkgurl

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatchesScpLike(t *testing.T) {
	// See https://github.com/git/git/blob/master/Documentation/urls.txt#L37
	examples := []string{
		// Most-extended case
		"git@github.com:james/bond",
		// Most-extended case with port
		"git@github.com:22:james/bond",
		// Most-extended case with numeric path
		"git@github.com:007/bond",
		// Most-extended case with port and numeric "username"
		"git@github.com:22:007/bond",
		// Single repo path
		"git@github.com:bond",
		// Single repo path with port
		"git@github.com:22:bond",
		// Single repo path with port and numeric repo
		"git@github.com:22:007",
		// Repo path ending with .git and starting with _
		"git@github.com:22:_007.git",
		"git@github.com:_007.git",
		"git@github.com:_james.git",
		"git@github.com:_james/bond.git",
		"git@github.com:_james/bond.git#nested/subdir",
	}

	for _, url := range examples {
		t.Run(url, func(t *testing.T) {
			assert.Equal(t, true, MatchesScpLike(url))
		})
	}
}

func TestFindScpLikeComponents(t *testing.T) {
	testCases := []struct {
		url, user, host, port, path, subdir string
	}{
		{
			// Most-extended case
			url: "git@github.com:james/bond", user: "git", host: "github.com", port: "", path: "james/bond", subdir: "",
		},
		{
			// Most-extended case with port
			url: "git@github.com:22:james/bond", user: "git", host: "github.com", port: "22", path: "james/bond", subdir: "",
		},
		{
			// Most-extended case with numeric path
			url: "git@github.com:007/bond", user: "git", host: "github.com", port: "", path: "007/bond", subdir: "",
		},
		{
			// Most-extended case with port and numeric path
			url: "git@github.com:22:007/bond", user: "git", host: "github.com", port: "22", path: "007/bond", subdir: "",
		},
		{
			// Single repo path
			url: "git@github.com:bond", user: "git", host: "github.com", port: "", path: "bond", subdir: "",
		},
		{
			// Single repo path with subdirectory
			url: "git@github.com:bond#subdir", user: "git", host: "github.com", port: "", path: "bond", subdir: "subdir",
		},
		{
			// Single repo path with port
			url: "git@github.com:22:bond", user: "git", host: "github.com", port: "22", path: "bond", subdir: "",
		},
		{
			// Single repo path with port and numeric path
			url: "git@github.com:22:007", user: "git", host: "github.com", port: "22", path: "007", subdir: "",
		},
		{
			// Repo path ending with .git and starting with _
			url: "git@github.com:22:_007.git", user: "git", host: "github.com", port: "22", path: "_007.git", subdir: "",
		},
		{
			// Repo path ending with .git and starting with _
			url: "git@github.com:_007.git", user: "git", host: "github.com", port: "", path: "_007.git", subdir: "",
		},
		{
			// Repo path ending with .git and starting with _
			url: "git@github.com:_james.git", user: "git", host: "github.com", port: "", path: "_james.git", subdir: "",
		},
		{
			// Repo path ending with .git and starting with _
			url: "git@github.com:_james/bond.git", user: "git", host: "github.com", port: "", path: "_james/bond.git", subdir: "",
		},
	}

	for _, tc := range testCases {
		user, host, port, path, subdir := FindScpLikeComponents(tc.url)

		assert.Equal(t, tc.user, user)
		assert.Equal(t, tc.host, host)
		assert.Equal(t, tc.port, port)
		assert.Equal(t, tc.path, path)
		assert.Equal(t, tc.subdir, subdir)
	}
}
