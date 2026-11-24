package commands

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/scaffoldrc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const readmeSnippet = `---
questions:
  - name: description
    prompt:
      message: Description
presets:
  default:
    description: from preset
computed:
  shout: "{{ upper .Scaffold.description }}"
---

{{ .Scaffold.description }} {{ .Computed.shout }} {{ .Scaffold.author }}
`

func newSnipController(t *testing.T, rc *scaffoldrc.ScaffoldRC) (*Controller, string) {
	t.Helper()

	rcDir := t.TempDir()
	snipDir := filepath.Join(rcDir, "snippets")
	require.NoError(t, os.MkdirAll(snipDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(snipDir, "readme.md"), []byte(readmeSnippet), 0o644))

	ctrl := &Controller{Flags: Flags{ScaffoldRCPath: filepath.Join(rcDir, "scaffoldrc.yml")}}
	ctrl.Prepare(engine.New(), rc)

	return ctrl, rcDir
}

func TestSnip_NoPromptWritesFile(t *testing.T) {
	rc := scaffoldrc.Default()
	rc.Defaults = map[string]any{"author": "rc-author", "description": "from rc"}
	ctrl, _ := newSnipController(t, rc)

	out := t.TempDir()
	err := ctrl.Snip([]string{"readme", "description=from cli"}, FlagsSnip{NoPrompt: true, Preset: "default", Output: out})
	require.NoError(t, err)

	got, err := os.ReadFile(filepath.Join(out, "readme.md"))
	require.NoError(t, err)
	assert.Equal(t, "from cli FROM CLI rc-author\n", string(got))
}

func TestSnip_NoPromptRefusesToClobber(t *testing.T) {
	ctrl, _ := newSnipController(t, scaffoldrc.Default())

	out := t.TempDir()
	dest := filepath.Join(out, "readme.md")
	require.NoError(t, os.WriteFile(dest, []byte("keep"), 0o644))

	err := ctrl.Snip([]string{"readme.md"}, FlagsSnip{NoPrompt: true, Preset: "default", Output: out})
	require.ErrorContains(t, err, "already exists")

	got, err := os.ReadFile(dest)
	require.NoError(t, err)
	assert.Equal(t, "keep", string(got))
}

func TestSnip_NoPromptErrors(t *testing.T) {
	ctrl, _ := newSnipController(t, scaffoldrc.Default())

	err := ctrl.Snip(nil, FlagsSnip{NoPrompt: true})
	require.ErrorContains(t, err, "snippet name is required")

	err = ctrl.Snip([]string{"missing"}, FlagsSnip{NoPrompt: true})
	require.ErrorContains(t, err, "not found")

	err = ctrl.Snip([]string{"readme"}, FlagsSnip{NoPrompt: true, Preset: "nope"})
	require.ErrorContains(t, err, "preset 'nope' not found")
}

func TestSnippetRoots(t *testing.T) {
	home, err := os.UserHomeDir()
	require.NoError(t, err)

	tests := []struct {
		name     string
		snippets []string
		flags    []string
		want     []string
	}{
		{
			name: "default next to scaffoldrc",
			want: []string{"/rc/snippets"},
		},
		{
			name:     "configured entries replace the default",
			snippets: []string{"~/dots/snips", "local", "/abs/**/*.md"},
			want:     []string{filepath.Join(home, "dots/snips"), "/rc/local", "/abs/**/*.md"},
		},
		{
			name:  "flag entries come after the default",
			flags: []string{"extra"},
			want:  []string{"/rc/snippets", "extra"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := scaffoldrc.Default()
			rc.Snippets = tt.snippets

			ctrl := &Controller{Flags: Flags{ScaffoldRCPath: "/rc/scaffoldrc.yml", SnippetDirs: tt.flags}}
			ctrl.Prepare(engine.New(), rc)

			assert.Equal(t, tt.want, ctrl.snippetRoots())
		})
	}
}

func TestWriteSnippet(t *testing.T) {
	t.Run("directory destination uses the snippet file name", func(t *testing.T) {
		dir := t.TempDir()

		require.NoError(t, writeSnippet(dir, "a.txt", []byte("a"), 0o644, false))
		assert.FileExists(t, filepath.Join(dir, "a.txt"))
	})

	t.Run("file destination creates parent directories", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "nested", "b.txt")

		require.NoError(t, writeSnippet(dest, "a.txt", []byte("b"), 0o644, false))

		info, err := os.Stat(dest)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o644), info.Mode().Perm())
	})

	t.Run("trailing separator creates the directory", func(t *testing.T) {
		dir := filepath.Join(t.TempDir(), "new") + string(filepath.Separator)

		require.NoError(t, writeSnippet(dir, "d.txt", []byte("d"), 0o644, false))
		assert.FileExists(t, filepath.Join(dir, "d.txt"))
	})

	t.Run("perm is applied", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "run.sh")

		require.NoError(t, writeSnippet(dest, "run.sh", []byte("#!/bin/sh"), 0o755, false))

		info, err := os.Stat(dest)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
	})

	t.Run("overwrite applies perm to an existing file", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "run.sh")
		require.NoError(t, os.WriteFile(dest, []byte("old"), 0o644))

		require.NoError(t, writeSnippet(dest, "run.sh", []byte("#!/bin/sh"), 0o755, true))

		info, err := os.Stat(dest)
		require.NoError(t, err)
		assert.Equal(t, os.FileMode(0o755), info.Mode().Perm())
	})

	t.Run("overwrite replaces existing content", func(t *testing.T) {
		dest := filepath.Join(t.TempDir(), "c.txt")
		require.NoError(t, os.WriteFile(dest, []byte("old content"), 0o644))

		require.NoError(t, writeSnippet(dest, "c.txt", []byte("new"), 0o644, true))

		got, err := os.ReadFile(dest)
		require.NoError(t, err)
		assert.Equal(t, "new", string(got))
	})
}

func TestSnippetPerm(t *testing.T) {
	assert.Equal(t, os.FileMode(0o644), snippetPerm(0o664))
	assert.Equal(t, os.FileMode(0o755), snippetPerm(0o775))
	assert.Equal(t, os.FileMode(0o755), snippetPerm(0o700))
}

func TestSetSnippetProject(t *testing.T) {
	p := &scaffold.Project{Name: scaffold.TemplateDirName}
	vars := map[string]any{"Project": scaffold.TemplateDirName}
	setSnippetProject(p, vars, map[string]any{"Project": "myapp"})
	assert.Equal(t, "myapp", p.Name)
	assert.Equal(t, "myapp", vars["Project"])

	p = &scaffold.Project{Name: scaffold.TemplateDirName}
	vars = map[string]any{"Project": scaffold.TemplateDirName}
	setSnippetProject(p, vars, map[string]any{})
	assert.Empty(t, p.Name)
	assert.NotContains(t, vars, "Project")
}
