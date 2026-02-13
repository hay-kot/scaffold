package scaffold

import (
	"testing"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v3"
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
		"ProjectSlug",
		"ProjectCamel",
		"ProjectPascal",
	}

	for _, key := range requiredStringKeys {
		assert.NotNil(t, got[key])
		assert.IsType(t, "", got[key])
	}
	assert.Equal(t, got["ProjectKebab"], got["ProjectSlug"])

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

func Test_detectEachPattern(t *testing.T) {
	tests := []struct {
		path      string
		wantVar   string
		wantToken string
	}{
		{"{{ .Project }}/[services]/handler.go", "services", "[services]"},
		{"{{ .Project }}/[models].go", "models", "[models]"},
		{"{{ .Project }}/normal/file.go", "", ""},
		{"[items].txt", "items", "[items]"},
		{"[_private]/file.go", "_private", "[_private]"},
		{"path/[var1]/[var2]/file.go", "var1", "[var1]"},
		{"no-brackets/here.txt", "", ""},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			gotVar, gotToken := detectEachPattern(tt.path)
			assert.Equal(t, tt.wantVar, gotVar)
			assert.Equal(t, tt.wantToken, gotToken)
		})
	}
}

func Test_resolveListVar(t *testing.T) {
	tests := []struct {
		name    string
		vars    engine.Vars
		varName string
		want    []string
		wantErr bool
	}{
		{
			name: "string slice",
			vars: engine.Vars{
				"Scaffold": engine.Vars{
					"services": []string{"auth", "users"},
				},
			},
			varName: "services",
			want:    []string{"auth", "users"},
		},
		{
			name: "any slice of strings",
			vars: engine.Vars{
				"Scaffold": engine.Vars{
					"items": []any{"foo", "bar"},
				},
			},
			varName: "items",
			want:    []string{"foo", "bar"},
		},
		{
			name: "plain string wraps to single-element slice",
			vars: engine.Vars{
				"Scaffold": engine.Vars{
					"name": "single",
				},
			},
			varName: "name",
			want:    []string{"single"},
		},
		{
			name: "missing variable",
			vars: engine.Vars{
				"Scaffold": engine.Vars{},
			},
			varName: "missing",
			wantErr: true,
		},
		{
			name: "wrong type",
			vars: engine.Vars{
				"Scaffold": engine.Vars{
					"count": 42,
				},
			},
			varName: "count",
			wantErr: true,
		},
		{
			name:    "no scaffold key",
			vars:    engine.Vars{},
			varName: "anything",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolveListVar(tt.vars, tt.varName)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_EachConfig_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		name string
		yaml string
		want []EachConfig
	}{
		{
			name: "string shorthand",
			yaml: "each:\n  - services\n  - models\n",
			want: []EachConfig{
				{Var: "services"},
				{Var: "models"},
			},
		},
		{
			name: "object form",
			yaml: "each:\n  - var: models\n    as: \"{{ .Each.Item | toPascalCase }}\"\n",
			want: []EachConfig{
				{Var: "models", As: "{{ .Each.Item | toPascalCase }}"},
			},
		},
		{
			name: "mixed",
			yaml: "each:\n  - services\n  - var: models\n    as: \"{{ .Each.Item | toPascalCase }}\"\n",
			want: []EachConfig{
				{Var: "services"},
				{Var: "models", As: "{{ .Each.Item | toPascalCase }}"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var conf ProjectScaffoldFile
			err := yaml.Unmarshal([]byte(tt.yaml), &conf)
			require.NoError(t, err)
			assert.Equal(t, tt.want, conf.Each)
		})
	}
}

func Test_isEachVar(t *testing.T) {
	configs := []EachConfig{
		{Var: "services"},
		{Var: "models", As: "{{ .Each.Item | toPascalCase }}"},
	}

	ec, ok := isEachVar(configs, "services")
	assert.True(t, ok)
	assert.Equal(t, "services", ec.Var)
	assert.Empty(t, ec.As)

	ec, ok = isEachVar(configs, "models")
	assert.True(t, ok)
	assert.Equal(t, "{{ .Each.Item | toPascalCase }}", ec.As)

	_, ok = isEachVar(configs, "notfound")
	assert.False(t, ok)
}
