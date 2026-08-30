package scaffold

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ComputedOrder_FollowsDeclaration(t *testing.T) {
	conf, err := ReadScaffoldFile(strings.NewReader(`
computed:
  zeta: "1"
  alpha: "2"
  middle: "3"
`))
	require.NoError(t, err)

	assert.Equal(t, map[string]string{"zeta": "1", "alpha": "2", "middle": "3"}, conf.Computed,
		"Computed stays a plain map")
	assert.Equal(t, []string{"zeta", "alpha", "middle"}, conf.ComputedOrder(),
		"order must follow scaffold.yaml, not map iteration")
}

func Test_ComputedOrder_Empty(t *testing.T) {
	conf, err := ReadScaffoldFile(strings.NewReader("questions: []\n"))
	require.NoError(t, err)

	assert.Nil(t, conf.ComputedOrder())
}

// Test_ComputedOrder_ProgrammaticIsSorted covers a caller that builds the struct
// directly rather than reading YAML. There is no declared order to follow, so
// sorting keeps a render reproducible instead of leaving it to map iteration.
func Test_ComputedOrder_ProgrammaticIsSorted(t *testing.T) {
	conf := &ProjectScaffoldFile{
		Computed: map[string]string{"zeta": "1", "alpha": "2", "middle": "3"},
	}

	assert.Equal(t, []string{"alpha", "middle", "zeta"}, conf.ComputedOrder())
}

// Test_ComputedOrder_MixedSources covers a caller that reads YAML and then adds a
// key. Declared keys keep their order and the added key follows.
func Test_ComputedOrder_MixedSources(t *testing.T) {
	conf, err := ReadScaffoldFile(strings.NewReader(`
computed:
  zeta: "1"
  alpha: "2"
`))
	require.NoError(t, err)

	conf.Computed["added"] = "3"

	assert.Equal(t, []string{"zeta", "alpha", "added"}, conf.ComputedOrder())
}

// Test_ComputedOrder_DroppedKey covers a caller that removes a declared key. The
// stale entry in the recorded order must not produce a phantom key.
func Test_ComputedOrder_DroppedKey(t *testing.T) {
	conf, err := ReadScaffoldFile(strings.NewReader(`
computed:
  zeta: "1"
  alpha: "2"
`))
	require.NoError(t, err)

	delete(conf.Computed, "zeta")

	assert.Equal(t, []string{"alpha"}, conf.ComputedOrder())
}

// Test_BuildVars_ComputedChaining is the reason order is recorded: a port
// scaffold derives a block once and offsets from it, rather than repeating the
// seed expression for every service.
func Test_BuildVars_ComputedChaining(t *testing.T) {
	conf, err := ReadScaffoldFile(strings.NewReader(`
computed:
  base: "{{ portBlock .Project }}"
  api: "{{ add .Computed.base 1 }}"
  web: "{{ add .Computed.base 2 }}"
  label: "{{ .Project }}-{{ .Computed.api }}"
`))
	require.NoError(t, err)

	project := &Project{Name: "recipinned", Conf: conf}

	got, err := BuildVars(tEngine, project, engine.Vars{})
	require.NoError(t, err)

	computed, ok := got["Computed"].(map[string]any)
	require.True(t, ok)

	base, ok := computed["base"].(int)
	require.True(t, ok, "portBlock output must coerce to int")

	assert.Equal(t, base+1, computed["api"])
	assert.Equal(t, base+2, computed["web"])
	assert.Equal(t, "recipinned-"+strconv.Itoa(base+1), computed["label"])
}

func Test_BuildVars_ComputedForwardReferenceIsEmpty(t *testing.T) {
	// Referencing a later declaration is not an error, it just has no value yet.
	// Keeping this explicit so the ordering contract is visible in the tests.
	conf, err := ReadScaffoldFile(strings.NewReader(`
computed:
  early: "{{ .Computed.late }}"
  late: "set"
`))
	require.NoError(t, err)

	got, err := BuildVars(tEngine, &Project{Name: "demo", Conf: conf}, engine.Vars{})
	require.NoError(t, err)

	computed := got["Computed"].(map[string]any)
	assert.Equal(t, "<no value>", computed["early"])
	assert.Equal(t, "set", computed["late"])
}

func Test_BuildVars_ComputedErrorNamesTheKey(t *testing.T) {
	project := &Project{
		Name: "demo",
		Conf: &ProjectScaffoldFile{
			Computed: map[string]string{"broken": `{{ portBlock "" }}`},
		},
	}

	_, err := BuildVars(tEngine, project, engine.Vars{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), `computed "broken"`)
}

func Test_BuildVars_Ctx(t *testing.T) {
	dir := t.TempDir()

	project := &Project{
		Name:    "demo",
		Conf:    &ProjectScaffoldFile{},
		Options: Options{OutputDir: dir},
	}

	got, err := BuildVars(tEngine, project, engine.Vars{})
	require.NoError(t, err)

	ctx, ok := got["Ctx"].(engine.Vars)
	require.True(t, ok)

	assert.Equal(t, dir, ctx["OutputDir"])
	assert.Equal(t, filepath.Base(dir), ctx["OutputDirBase"])

	// t.TempDir on macOS hands back a /var symlink, so compare against the same
	// resolution filepath.Abs performs rather than the raw string.
	wantAbs, err := filepath.Abs(dir)
	require.NoError(t, err)
	assert.Equal(t, wantAbs, ctx["OutputDirAbs"])

	wd, err := os.Getwd()
	require.NoError(t, err)
	assert.Equal(t, wd, ctx["WorkingDir"])
	assert.Equal(t, filepath.Base(wd), ctx["WorkingDirBase"])

	gitVars, ok := ctx["Git"].(engine.Vars)
	require.True(t, ok)
	assert.Contains(t, gitVars, "Branch")
	assert.Contains(t, gitVars, "Repo")
}

func Test_BuildVars_CtxMemoryOutputDir(t *testing.T) {
	project := &Project{
		Name:    "demo",
		Conf:    &ProjectScaffoldFile{},
		Options: Options{OutputDir: memoryOutputDir},
	}

	got, err := BuildVars(tEngine, project, engine.Vars{})
	require.NoError(t, err)

	ctx := got["Ctx"].(engine.Vars)
	assert.Equal(t, memoryOutputDir, ctx["OutputDir"])
	assert.Empty(t, ctx["OutputDirAbs"], "the sentinel is not a real path")
	assert.Empty(t, ctx["OutputDirBase"])
}

func Test_BuildVars_CtxEmptyOutputDir(t *testing.T) {
	project := &Project{
		Name: "demo",
		Conf: &ProjectScaffoldFile{},
	}

	got, err := BuildVars(tEngine, project, engine.Vars{})
	require.NoError(t, err)

	ctx := got["Ctx"].(engine.Vars)
	assert.Empty(t, ctx["OutputDir"])
	assert.Empty(t, ctx["OutputDirBase"])
	assert.NotEmpty(t, ctx["WorkingDir"], "working dir stays available as a fallback seed")
}
