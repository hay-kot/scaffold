package scaffold

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rwfs"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func readAll(t *testing.T, fsys fs.FS, path string) string {
	t.Helper()

	bits, err := fs.ReadFile(fsys, path)
	require.NoError(t, err)
	return string(bits)
}

// Test_NoClobber_TemplateScaffoldChecksOutputPath covers the case that made
// no-clobber useless in a project: guards run against "templates/.env" while the
// file on disk is ".env", so the guard never matched and overwrote local edits.
func Test_NoClobber_TemplateScaffoldChecksOutputPath(t *testing.T) {
	readFS := fstest.MapFS{
		"scaffold.yaml":     &fstest.MapFile{Data: []byte("questions: []\n")},
		"templates/.env":    &fstest.MapFile{Data: []byte("GENERATED=true\n")},
		"templates/new.txt": &fstest.MapFile{Data: []byte("fresh\n")},
	}

	memFS := rwfs.NewMemoryWFS()
	require.NoError(t, memFS.WriteFile(".env", []byte("HAND_EDITED=keepme\n"), 0o644))

	project := &Project{
		NameTemplate: TemplateDirName,
		Name:         TemplateDirName,
		Conf:         &ProjectScaffoldFile{},
		Options:      Options{NoClobber: true},
	}

	args := &RWFSArgs{ReadFS: readFS, WriteFS: memFS, Project: project}

	vars, err := BuildVars(tEngine, project, engine.Vars{})
	require.NoError(t, err)

	require.NoError(t, RenderRWFS(tEngine, args, vars))

	assert.Equal(t, "HAND_EDITED=keepme\n", readAll(t, memFS, ".env"),
		"existing file must survive a re-run")
	assert.Equal(t, "fresh\n", readAll(t, memFS, "new.txt"),
		"a collision must not stop the files that follow it")
}

// Test_NoClobber_ContinuesPastCollision covers the second half of the bug: the
// guard returned errFileExists into the walk, which aborted the whole render.
func Test_NoClobber_ContinuesPastCollision(t *testing.T) {
	readFS := fstest.MapFS{
		"scaffold.yaml":                    &fstest.MapFile{Data: []byte("questions: []\n")},
		"{{ .Project }}/a.txt":             &fstest.MapFile{Data: []byte("a from template\n")},
		"{{ .Project }}/b.txt":             &fstest.MapFile{Data: []byte("b from template\n")},
		"{{ .Project }}/nested/c.txt":      &fstest.MapFile{Data: []byte("c from template\n")},
		"{{ .Project }}/nested/deep/d.txt": &fstest.MapFile{Data: []byte("d from template\n")},
	}

	memFS := rwfs.NewMemoryWFS()
	require.NoError(t, memFS.MkdirAll("demo", dirPerm))
	require.NoError(t, memFS.WriteFile("demo/a.txt", []byte("edited\n"), 0o644))

	project := &Project{
		NameTemplate: "{{ .Project }}",
		Name:         "demo",
		Conf:         &ProjectScaffoldFile{},
		Options:      Options{NoClobber: true},
	}

	args := &RWFSArgs{ReadFS: readFS, WriteFS: memFS, Project: project}

	vars, err := BuildVars(tEngine, project, engine.Vars{})
	require.NoError(t, err)

	require.NoError(t, RenderRWFS(tEngine, args, vars))

	assert.Equal(t, "edited\n", readAll(t, memFS, "demo/a.txt"))
	assert.Equal(t, "b from template\n", readAll(t, memFS, "demo/b.txt"))
	assert.Equal(t, "c from template\n", readAll(t, memFS, "demo/nested/c.txt"))
	assert.Equal(t, "d from template\n", readAll(t, memFS, "demo/nested/deep/d.txt"))
}

func Test_NoClobber_DisabledOverwrites(t *testing.T) {
	readFS := fstest.MapFS{
		"scaffold.yaml":  &fstest.MapFile{Data: []byte("questions: []\n")},
		"templates/.env": &fstest.MapFile{Data: []byte("GENERATED=true\n")},
	}

	memFS := rwfs.NewMemoryWFS()
	require.NoError(t, memFS.WriteFile(".env", []byte("HAND_EDITED=keepme\n"), 0o644))

	project := &Project{
		NameTemplate: TemplateDirName,
		Name:         TemplateDirName,
		Conf:         &ProjectScaffoldFile{},
		Options:      Options{NoClobber: false},
	}

	args := &RWFSArgs{ReadFS: readFS, WriteFS: memFS, Project: project}

	vars, err := BuildVars(tEngine, project, engine.Vars{})
	require.NoError(t, err)

	require.NoError(t, RenderRWFS(tEngine, args, vars))

	assert.Equal(t, "GENERATED=true\n", readAll(t, memFS, ".env"),
		"--overwrite must still replace the file")
}

// Test_NoClobber_EachExpansion checks the expanded-file paths, which run through a
// separate guard loop from the normal walk.
func Test_NoClobber_EachExpansion(t *testing.T) {
	readFS := fstest.MapFS{
		"scaffold.yaml":               &fstest.MapFile{Data: []byte("questions: []\n")},
		"templates/[services].conf":   &fstest.MapFile{Data: []byte("service={{ .Each.Item }}\n")},
		"templates/[services]/run.sh": &fstest.MapFile{Data: []byte("# {{ .Each.Item }}\n")},
	}

	memFS := rwfs.NewMemoryWFS()
	require.NoError(t, memFS.WriteFile("api.conf", []byte("hand written\n"), 0o644))

	project := &Project{
		NameTemplate: TemplateDirName,
		Name:         TemplateDirName,
		Conf: &ProjectScaffoldFile{
			Each: []EachConfig{{Var: "services"}},
		},
		Options: Options{NoClobber: true},
	}

	args := &RWFSArgs{ReadFS: readFS, WriteFS: memFS, Project: project}

	vars, err := BuildVars(tEngine, project, engine.Vars{"services": []string{"api", "web"}})
	require.NoError(t, err)

	require.NoError(t, RenderRWFS(tEngine, args, vars))

	assert.Equal(t, "hand written\n", readAll(t, memFS, "api.conf"))
	assert.Equal(t, "service=web\n", readAll(t, memFS, "web.conf"))
	assert.Equal(t, "# api\n", readAll(t, memFS, "api/run.sh"))
	assert.Equal(t, "# web\n", readAll(t, memFS, "web/run.sh"))
}

func Test_FilePerm(t *testing.T) {
	tests := []struct {
		name string
		mode fs.FileMode
		want fs.FileMode
	}{
		{name: "plain config stays 0644", mode: 0o644, want: 0o644},
		{name: "executable template stays executable", mode: 0o755, want: 0o755},
		{name: "read only source gains owner write", mode: 0o444, want: 0o644},
		{name: "group and world write are cleared", mode: 0o666, want: 0o644},
		{name: "world writable executable is narrowed", mode: 0o777, want: 0o755},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsys := fstest.MapFS{"f": &fstest.MapFile{Data: []byte("x"), Mode: tt.mode}}

			entries, err := fs.ReadDir(fsys, ".")
			require.NoError(t, err)
			require.Len(t, entries, 1)

			assert.Equal(t, tt.want, filePerm(entries[0]))
		})
	}
}

func Test_FilePerm_NilEntry(t *testing.T) {
	assert.Equal(t, defaultFilePerm, filePerm(nil))
}
