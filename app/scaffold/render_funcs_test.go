package scaffold

import (
	"testing"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_BuildVars(t *testing.T) {
	project := &Project{
		Name: "test",
		Conf: &ProjectScaffoldFile{
			Computed: map[string]string{
				"Bool":     "true",
				"Int":      "{{ add 1 2 }}",
				"ZeroInt":  "0",
				"UsesVars": "Computed: {{ .Scaffold.key }}",
			},
		},
	}

	vars := engine.Vars{
		"key": "value",
	}

	eng := engine.New()

	got, err := BuildVars(eng, project, vars)
	require.NoError(t, err)

	//
	// Assert Top Level Keys
	//

	requiredStringKeys := []string{
		"Project",
		"ProjectSnake",
		"ProjectKebab",
		"ProjectCamel",
	}

	for _, key := range requiredStringKeys {
		assert.NotNil(t, got[key])
		assert.IsType(t, "", got[key])
	}

	assert.NotNil(t, got["Scaffold"])
	assert.IsType(t, engine.Vars{}, got["Scaffold"])

	//
	// Assert Passed in Vars live under Scaffold
	//

	scaffold := got["Scaffold"].(engine.Vars)
	assert.NotNil(t, scaffold["key"])

	//
	// Assert Computed Properties are Typed and Computed
	//

	require.NotNil(t, got["Computed"])
	computed := got["Computed"].(map[string]any)

	assert.NotNil(t, computed["Bool"])
	assert.IsType(t, true, computed["Bool"])

	assert.NotNil(t, computed["Int"])
	assert.IsType(t, 3, computed["Int"])

	assert.NotNil(t, computed["ZeroInt"])
	assert.IsType(t, 0, computed["ZeroInt"])

	assert.NotNil(t, computed["UsesVars"])
	assert.IsType(t, "", computed["UsesVars"])
	assert.Equal(t, "Computed: value", computed["UsesVars"])
}
