package snippets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeFiles(t *testing.T, root string, files ...string) {
	t.Helper()
	for _, f := range files {
		p := filepath.Join(root, filepath.FromSlash(f))
		require.NoError(t, os.MkdirAll(filepath.Dir(p), 0o755))
		require.NoError(t, os.WriteFile(p, []byte(f), 0o644))
	}
}

func TestDiscover_Directory(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, "taskfile.yml", "go/main.go", ".hidden", ".git/config", "go/.DS_Store")

	got, err := Discover([]string{root})
	require.NoError(t, err)

	assert.Equal(t, []string{"go/main.go", "taskfile.yml"}, Names(got))
	assert.Equal(t, filepath.Join(root, "go", "main.go"), got[0].Path)
}

func TestDiscover_Glob(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, "a.md", "nested/b.md", "nested/c.txt")

	got, err := Discover([]string{filepath.Join(root, "**", "*.md")})
	require.NoError(t, err)

	assert.Equal(t, []string{"a.md", "nested/b.md"}, Names(got))
	assert.Equal(t, filepath.Join(root, "nested", "b.md"), got[1].Path)
}

func TestDiscover_SingleFile(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, "one.txt")

	got, err := Discover([]string{filepath.Join(root, "one.txt")})
	require.NoError(t, err)

	assert.Equal(t, []string{"one.txt"}, Names(got))
}

func TestDiscover_MissingRootsSkipped(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	got, err := Discover([]string{missing, filepath.Join(missing, "*.md")})
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestDiscover_FirstRootWins(t *testing.T) {
	first, second := t.TempDir(), t.TempDir()
	writeFiles(t, first, "dup.txt")
	writeFiles(t, second, "dup.txt", "other.txt")

	got, err := Discover([]string{first, second})
	require.NoError(t, err)

	assert.Equal(t, []string{"dup.txt", "other.txt"}, Names(got))
	assert.Equal(t, filepath.Join(first, "dup.txt"), got[0].Path)
}

func TestFind(t *testing.T) {
	list := []Snippet{
		{Name: "readme.md"},
		{Name: "readme"},
		{Name: "go/main.go"},
	}

	tests := []struct {
		name string
		want string
		ok   bool
	}{
		{name: "readme.md", want: "readme.md", ok: true},
		{name: "readme", want: "readme", ok: true},
		{name: "go/main", want: "go/main.go", ok: true},
		{name: "main", ok: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := Find(list, tt.name)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got.Name)
		})
	}
}

func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)
	t.Setenv("SNIPPET_TEST_DIR", "/env/dir")

	assert.Equal(t, filepath.Join(home, "snips"), ExpandPath("~/snips"))
	assert.Equal(t, "/env/dir/x", ExpandPath("$SNIPPET_TEST_DIR/x"))
	assert.Equal(t, "/abs", ExpandPath("/abs"))
}

func TestParse(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantBody  string
		wantQuest int
		wantPre   string
		wantErr   bool
	}{
		{
			name: "front matter with blank separator",
			input: `---
messages:
  pre: hello
questions:
  - name: description
    prompt:
      message: Description
---

Description: {{ .Scaffold.description }}
`,
			wantBody:  "Description: {{ .Scaffold.description }}\n",
			wantQuest: 1,
			wantPre:   "hello",
		},
		{
			name:     "crlf line endings",
			input:    "---\r\nmessages:\r\n  pre: hi\r\n---\r\nbody\r\n",
			wantBody: "body\r\n",
			wantPre:  "hi",
		},
		{
			name:     "no front matter",
			input:    "just a template {{ .Project }}\n",
			wantBody: "just a template {{ .Project }}\n",
		},
		{
			name:     "empty front matter",
			input:    "---\n---\nbody",
			wantBody: "body",
		},
		{
			name:     "only one blank line dropped",
			input:    "---\n---\n\n\nbody",
			wantBody: "\nbody",
		},
		{
			name:    "unterminated front matter",
			input:   "---\nmessages:\n  pre: hi\n",
			wantErr: true,
		},
		{
			name:    "invalid yaml",
			input:   "---\nquestions: [\n---\nbody",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			conf, body, err := Parse(strings.NewReader(tt.input))
			if tt.wantErr {
				require.Error(t, err)
				return
			}

			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, body)
			assert.Len(t, conf.Questions, tt.wantQuest)
			assert.Equal(t, tt.wantPre, conf.Messages.Pre)
		})
	}
}
