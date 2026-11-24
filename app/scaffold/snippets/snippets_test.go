package snippets

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/hay-kot/scaffold/app/core/apperrors"
	"github.com/hay-kot/scaffold/app/core/engine"
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

	got := Discover([]string{root})

	assert.Equal(t, []string{"go/main.go", "taskfile.yml"}, Names(got))
	assert.Equal(t, filepath.Join(root, "go", "main.go"), got[0].Path)
}

func TestDiscover_Glob(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, "a.md", "nested/b.md", "nested/c.txt")

	got := Discover([]string{filepath.Join(root, "**", "*.md")})

	assert.Equal(t, []string{"a.md", "nested/b.md"}, Names(got))
	assert.Equal(t, filepath.Join(root, "nested", "b.md"), got[1].Path)
}

func TestDiscover_SingleFile(t *testing.T) {
	root := t.TempDir()
	writeFiles(t, root, "one.txt")

	got := Discover([]string{filepath.Join(root, "one.txt")})

	assert.Equal(t, []string{"one.txt"}, Names(got))
}

func TestDiscover_MissingRootsSkipped(t *testing.T) {
	missing := filepath.Join(t.TempDir(), "nope")

	got := Discover([]string{missing, filepath.Join(missing, "*.md")})
	assert.Empty(t, got)
}

func TestDiscover_FirstRootWins(t *testing.T) {
	first, second := t.TempDir(), t.TempDir()
	writeFiles(t, first, "dup.txt")
	writeFiles(t, second, "dup.txt", "other.txt")

	got := Discover([]string{first, second})

	assert.Equal(t, []string{"dup.txt", "other.txt"}, Names(got))
	assert.Equal(t, filepath.Join(first, "dup.txt"), got[0].Path)
}

func TestDiscover_FollowsSymlinks(t *testing.T) {
	root, target := t.TempDir(), t.TempDir()
	writeFiles(t, target, "main.go")
	writeFiles(t, root, "file-target.txt")
	require.NoError(t, os.Symlink(target, filepath.Join(root, "go")))
	require.NoError(t, os.Symlink(filepath.Join(root, "file-target.txt"), filepath.Join(root, "linked.txt")))
	require.NoError(t, os.Symlink(root, filepath.Join(target, "loop")))
	require.NoError(t, os.Symlink(filepath.Join(root, "missing"), filepath.Join(root, "broken")))

	got := Discover([]string{root})

	assert.Equal(t, []string{"file-target.txt", "go/main.go", "linked.txt"}, Names(got))
}

func TestDiscover_LiteralPathWithGlobCharacters(t *testing.T) {
	root := filepath.Join(t.TempDir(), "[client] snippets")
	writeFiles(t, root, "a.txt")

	got := Discover([]string{root})

	assert.Equal(t, []string{"a.txt"}, Names(got))
}

func TestDiscover_UnreadableRootSkipped(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can read any directory")
	}

	bad, good := t.TempDir(), t.TempDir()
	writeFiles(t, good, "ok.txt")
	require.NoError(t, os.Chmod(bad, 0o000))
	t.Cleanup(func() { _ = os.Chmod(bad, 0o755) })

	got := Discover([]string{bad, good})

	assert.Equal(t, []string{"ok.txt"}, Names(got))
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
			got, ok, err := Find(list, tt.name)
			require.NoError(t, err)
			assert.Equal(t, tt.ok, ok)
			assert.Equal(t, tt.want, got.Name)
		})
	}
}

func TestFind_AmbiguousStem(t *testing.T) {
	list := []Snippet{{Name: "readme.md"}, {Name: "readme.txt"}}

	_, ok, err := Find(list, "readme")
	require.ErrorIs(t, err, ErrAmbiguous)
	assert.False(t, ok)
	assert.ErrorContains(t, err, "readme.md, readme.txt")

	got, ok, err := Find(list, "readme.txt")
	require.NoError(t, err)
	assert.True(t, ok)
	assert.Equal(t, "readme.txt", got.Name)
}

func TestDiscover_UnreadableSubdirSkipped(t *testing.T) {
	if os.Geteuid() == 0 {
		t.Skip("root can read any directory")
	}

	root := t.TempDir()
	writeFiles(t, root, "readme.md", "private/secret.txt")
	locked := filepath.Join(root, "private")
	require.NoError(t, os.Chmod(locked, 0o000))
	t.Cleanup(func() { _ = os.Chmod(locked, 0o755) })

	got := Discover([]string{root})

	assert.Equal(t, []string{"readme.md"}, Names(got))
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
		// BodyLine is 1 when the input has no front matter.
		wantBodyLine int
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
			wantBody:     "Description: {{ .Scaffold.description }}\n",
			wantQuest:    1,
			wantPre:      "hello",
			wantBodyLine: 10,
		},
		{
			name:         "crlf line endings",
			input:        "---\r\nmessages:\r\n  pre: hi\r\n---\r\nbody\r\n",
			wantBody:     "body\r\n",
			wantPre:      "hi",
			wantBodyLine: 5,
		},
		{
			name:         "no front matter",
			input:        "just a template {{ .Project }}\n",
			wantBody:     "just a template {{ .Project }}\n",
			wantBodyLine: 1,
		},
		{
			name:         "empty front matter",
			input:        "---\n---\nbody",
			wantBody:     "body",
			wantBodyLine: 3,
		},
		{
			name:         "only one blank line dropped",
			input:        "---\n---\n\n\nbody",
			wantBody:     "\nbody",
			wantBodyLine: 4,
		},
		{
			name:         "unterminated block is body",
			input:        "---\napiVersion: v1\nkind: Service\n",
			wantBody:     "---\napiVersion: v1\nkind: Service\n",
			wantBodyLine: 1,
		},
		{
			name:         "yaml document with unknown keys is body",
			input:        "---\napiVersion: v1\nkind: Service\n---\napiVersion: apps/v1\n",
			wantBody:     "---\napiVersion: v1\nkind: Service\n---\napiVersion: apps/v1\n",
			wantBodyLine: 1,
		},
		{
			name:         "markdown horizontal rules are body",
			input:        "---\nSome text\n---\nmore",
			wantBody:     "---\nSome text\n---\nmore",
			wantBodyLine: 1,
		},
		{
			name:         "invalid yaml is body",
			input:        "---\nquestions: [\n---\nbody",
			wantBody:     "---\nquestions: [\n---\nbody",
			wantBodyLine: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(strings.NewReader(tt.input))
			require.NoError(t, err)
			assert.Equal(t, tt.wantBody, got.Body)
			assert.Equal(t, tt.wantBodyLine, got.BodyLine)
			assert.Len(t, got.Conf.Questions, tt.wantQuest)
			assert.Equal(t, tt.wantPre, got.Conf.Messages.Pre)
		})
	}
}

func mustParse(t *testing.T, input string) *File {
	t.Helper()
	f, err := Parse(strings.NewReader(input))
	require.NoError(t, err)
	return f
}

func TestRender_Delimiters(t *testing.T) {
	const input = `---
delimiters:
  - glob: "workflows/*.yml"
    left: "[["
    right: "]]"
---
name: [[ .Scaffold.name ]]
run: ${{ github.sha }}
`

	t.Run("matching glob uses the override", func(t *testing.T) {
		got, err := mustParse(t, input).Render(engine.New(), "workflows/ci.yml", "ci.yml", engine.Vars{"Scaffold": map[string]any{"name": "ci"}})
		require.NoError(t, err)
		assert.Equal(t, "name: ci\nrun: ${{ github.sha }}\n", got)
	})

	t.Run("other names keep the default delimiters", func(t *testing.T) {
		got, err := mustParse(t, "---\ndelimiters:\n  - glob: \"*.yml\"\n    left: \"[[\"\n    right: \"]]\"\n---\n{{ .Scaffold.name }}").
			Render(engine.New(), "readme.md", "readme.md", engine.Vars{"Scaffold": map[string]any{"name": "x"}})
		require.NoError(t, err)
		assert.Equal(t, "x", got)
	})
}

func TestRender_EmptyBody(t *testing.T) {
	got, err := mustParse(t, "---\nmessages:\n  pre: hi\n---\n").Render(engine.New(), "empty", "empty", nil)
	require.NoError(t, err)
	assert.Empty(t, got)
}

func TestRender_ErrorReportsFileLine(t *testing.T) {
	const input = "---\nmessages:\n  pre: hi\n---\n\nline one\n{{ .Missing | nosuchfunc }}\nline three\n"

	_, err := mustParse(t, input).Render(engine.New(), "broken.txt", "/snips/broken.txt", nil)

	var terr *apperrors.TemplateError
	require.ErrorAs(t, err, &terr)
	assert.Equal(t, "/snips/broken.txt", terr.FilePath)
	assert.Equal(t, 7, terr.LineNumber)
	assert.Equal(t, []string{"line one", "{{ .Missing | nosuchfunc }}", "line three"}, terr.Context)
}

func TestIgnoredKeys(t *testing.T) {
	f := mustParse(t, `---
skip: ["a"]
rewrites:
  - from: a
    to: b
inject:
  - name: x
    path: y
    at: z
    template: t
features:
  - value: "true"
    globs: ["*"]
each: [items]
questions:
  - name: q
---
body`)

	assert.Equal(t, []string{"skip", "rewrites", "inject", "features", "each"}, f.IgnoredKeys())
	assert.Empty(t, mustParse(t, "body").IgnoredKeys())
}
