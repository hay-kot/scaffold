package scaffold

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rwfs"
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

func Test_ResolveFeatureVars(t *testing.T) {
	t.Run("nil vars", func(t *testing.T) {
		got := ResolveFeatureVars(nil)
		require.NotNil(t, got)
		require.NotNil(t, got["Scaffold"])
	})

	t.Run("unwraps Scaffold namespace into root scope", func(t *testing.T) {
		vars := engine.Vars{
			"Project": "test-project",
			"Scaffold": engine.Vars{
				"registry": "local",
				"builder":  "buildkit",
				"enabled":  true,
				"count":    5,
			},
			"Computed": map[string]any{
				"engine_type": "dockerd",
			},
		}

		resolved := ResolveFeatureVars(vars)

		// Bare variables available at root scope
		assert.Equal(t, "local", resolved["registry"])
		assert.Equal(t, "buildkit", resolved["builder"])
		assert.Equal(t, true, resolved["enabled"])
		assert.Equal(t, 5, resolved["count"])

		// Namespaced variables available under .Scaffold
		scaffoldMap, ok := resolved["Scaffold"].(engine.Vars)
		require.True(t, ok)
		assert.Equal(t, "local", scaffoldMap["registry"])
		assert.Equal(t, "buildkit", scaffoldMap["builder"])
		assert.Equal(t, true, scaffoldMap["enabled"])
		assert.Equal(t, 5, scaffoldMap["count"])

		// Root scope preserved
		assert.Equal(t, "test-project", resolved["Project"])
		assert.NotNil(t, resolved["Computed"])
	})

	t.Run("resolves bare root variables into Scaffold namespace", func(t *testing.T) {
		vars := engine.Vars{
			"registry": "local",
			"builder":  "buildkit",
			"enabled":  false,
		}

		resolved := ResolveFeatureVars(vars)

		// Bare variables available at root scope
		assert.Equal(t, "local", resolved["registry"])
		assert.Equal(t, "buildkit", resolved["builder"])
		assert.Equal(t, false, resolved["enabled"])

		// Resolved into .Scaffold
		scaffoldMap, ok := resolved["Scaffold"].(engine.Vars)
		require.True(t, ok)
		assert.Equal(t, "local", scaffoldMap["registry"])
		assert.Equal(t, "buildkit", scaffoldMap["builder"])
		assert.Equal(t, false, scaffoldMap["enabled"])
	})

	t.Run("resolves flat dotted Scaffold keys", func(t *testing.T) {
		vars := engine.Vars{
			"Scaffold.registry": "local",
			"Scaffold.builder":  "buildkit",
		}

		resolved := ResolveFeatureVars(vars)

		assert.Equal(t, "local", resolved["registry"])
		assert.Equal(t, "buildkit", resolved["builder"])

		scaffoldMap, ok := resolved["Scaffold"].(engine.Vars)
		require.True(t, ok)
		assert.Equal(t, "local", scaffoldMap["registry"])
		assert.Equal(t, "buildkit", scaffoldMap["builder"])
	})

	t.Run("preserves reserved root variables from being overwritten", func(t *testing.T) {
		computedMap := map[string]any{"type": "lima"}
		eachMap := map[string]any{"Item": "svc", "Index": 0}

		vars := engine.Vars{
			"Project":  "custom-app",
			"Computed": computedMap,
			"Each":     eachMap,
			"Scaffold": engine.Vars{
				"feature": "active",
			},
		}

		resolved := ResolveFeatureVars(vars)

		assert.Equal(t, computedMap, resolved["Computed"])
		assert.Equal(t, eachMap, resolved["Each"])
		assert.Equal(t, "custom-app", resolved["Project"])
		assert.Equal(t, "active", resolved["feature"])
	})
}

func Test_guardFeatureFlag_ScaffoldPrefixAndBareIdentical(t *testing.T) {
	eng := engine.New()

	testCases := []struct {
		name          string
		vars          engine.Vars
		prefixedValue string
		bareValue     string
		file          string
		pattern       string
		shouldInclude bool
	}{
		{
			name: "equality string match true",
			vars: engine.Vars{
				"Scaffold": engine.Vars{"registry": "local"},
			},
			prefixedValue: `{{ eq .Scaffold.registry "local" }}`,
			bareValue:     `{{ eq .registry "local" }}`,
			file:          "nested/registry-setup.yaml",
			pattern:       "**/registry*/**",
			shouldInclude: true,
		},
		{
			name: "equality string match false",
			vars: engine.Vars{
				"Scaffold": engine.Vars{"registry": "remote"},
			},
			prefixedValue: `{{ eq .Scaffold.registry "local" }}`,
			bareValue:     `{{ eq .registry "local" }}`,
			file:          "nested/registry-setup.yaml",
			pattern:       "**/registry*/**",
			shouldInclude: false,
		},
		{
			name: "bare vars input equality string match true",
			vars: engine.Vars{
				"registry": "local",
			},
			prefixedValue: `{{ eq .Scaffold.registry "local" }}`,
			bareValue:     `{{ eq .registry "local" }}`,
			file:          "nested/registry-setup.yaml",
			pattern:       "**/registry*/**",
			shouldInclude: true,
		},
		{
			name: "boolean flag true",
			vars: engine.Vars{
				"Scaffold": engine.Vars{"use_docker": true},
			},
			prefixedValue: `{{ .Scaffold.use_docker }}`,
			bareValue:     `{{ .use_docker }}`,
			file:          "docker/Dockerfile",
			pattern:       "docker/**",
			shouldInclude: true,
		},
		{
			name: "boolean flag false",
			vars: engine.Vars{
				"Scaffold": engine.Vars{"use_docker": false},
			},
			prefixedValue: `{{ .Scaffold.use_docker }}`,
			bareValue:     `{{ .use_docker }}`,
			file:          "docker/Dockerfile",
			pattern:       "docker/**",
			shouldInclude: false,
		},
		{
			name: "hasPrefix expression",
			vars: engine.Vars{
				"Scaffold": engine.Vars{"engine": "dockerd / compose"},
			},
			prefixedValue: `{{ hasPrefix "dockerd" .Scaffold.engine }}`,
			bareValue:     `{{ hasPrefix "dockerd" .engine }}`,
			file:          "compose/docker-compose.yaml",
			pattern:       "compose/**",
			shouldInclude: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Test with prefixed feature expression
			prefixedProject := &Project{
				Conf: &ProjectScaffoldFile{
					Features: []Feature{
						{
							Value: tc.prefixedValue,
							Globs: []string{tc.pattern},
						},
					},
				},
			}
			prefixedArgs := &RWFSArgs{Project: prefixedProject}
			prefixedGuard := guardFeatureFlag(eng, prefixedArgs, tc.vars)

			prefixedPath, prefixedErr := prefixedGuard(tc.file, nil)

			// Test with bare feature expression
			bareProject := &Project{
				Conf: &ProjectScaffoldFile{
					Features: []Feature{
						{
							Value: tc.bareValue,
							Globs: []string{tc.pattern},
						},
					},
				},
			}
			bareArgs := &RWFSArgs{Project: bareProject}
			bareGuard := guardFeatureFlag(eng, bareArgs, tc.vars)

			barePath, bareErr := bareGuard(tc.file, nil)

			// Both MUST evaluate identically
			assert.Equal(t, prefixedPath, barePath, "path results must match")
			assert.Equal(t, prefixedErr, bareErr, "error results must match")

			if tc.shouldInclude {
				assert.NoError(t, prefixedErr)
				assert.Equal(t, tc.file, prefixedPath)
			} else {
				assert.ErrorIs(t, prefixedErr, errSkipRender)
			}
		})
	}
}

func Test_RenderRWFS_FeatureFlagIdenticalEvaluation(t *testing.T) {
	eng := engine.New()

	scaffoldFiles := fstest.MapFS{
		"templates/common.txt":             &fstest.MapFile{Data: []byte("common content")},
		"templates/extra/registry-app.txt": &fstest.MapFile{Data: []byte("registry app content")},
	}

	tests := []struct {
		name          string
		vars          engine.Vars
		shouldInclude bool
	}{
		{
			name: "registry is local - feature enabled",
			vars: engine.Vars{
				"registry": "local",
			},
			shouldInclude: true,
		},
		{
			name: "registry is remote - feature disabled",
			vars: engine.Vars{
				"registry": "remote",
			},
			shouldInclude: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Run with .Scaffold prefix
			pPrefixed := &Project{
				Name:         "test-app",
				NameTemplate: "templates",
				Conf: &ProjectScaffoldFile{
					Features: []Feature{
						{
							Value: `{{ eq .Scaffold.registry "local" }}`,
							Globs: []string{"**/extra/**"},
						},
					},
				},
			}
			memFSPrefixed := rwfs.NewMemoryWFS()
			builtVarsPrefixed, err := BuildVars(eng, pPrefixed, tt.vars)
			require.NoError(t, err)

			argsPrefixed := &RWFSArgs{
				ReadFS:  scaffoldFiles,
				WriteFS: memFSPrefixed,
				Project: pPrefixed,
			}
			err = RenderRWFS(eng, argsPrefixed, builtVarsPrefixed)
			require.NoError(t, err)

			// Run with bare reference
			pBare := &Project{
				Name:         "test-app",
				NameTemplate: "templates",
				Conf: &ProjectScaffoldFile{
					Features: []Feature{
						{
							Value: `{{ eq .registry "local" }}`,
							Globs: []string{"**/extra/**"},
						},
					},
				},
			}
			memFSBare := rwfs.NewMemoryWFS()
			builtVarsBare, err := BuildVars(eng, pBare, tt.vars)
			require.NoError(t, err)

			argsBare := &RWFSArgs{
				ReadFS:  scaffoldFiles,
				WriteFS: memFSBare,
				Project: pBare,
			}
			err = RenderRWFS(eng, argsBare, builtVarsBare)
			require.NoError(t, err)

			// Verify identical files in both filesystems
			_, errPrefixed := memFSPrefixed.Open("extra/registry-app.txt")
			_, errBare := memFSBare.Open("extra/registry-app.txt")

			if tt.shouldInclude {
				assert.NoError(t, errPrefixed, "prefixed feature should include file")
				assert.NoError(t, errBare, "bare feature should include file")
			} else {
				assert.ErrorIs(t, errPrefixed, fs.ErrNotExist, "prefixed feature should exclude file")
				assert.ErrorIs(t, errBare, fs.ErrNotExist, "bare feature should exclude file")
			}
		})
	}
}

