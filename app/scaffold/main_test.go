package scaffold

import (
	"embed"
	"io/fs"
	"os"
	"testing"

	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	tEngine = engine.New()

	//go:embed testdata/projects/custom_delims
	customDelimsFiles embed.FS

	//go:embed testdata/projects/with_partials
	partialsFiles embed.FS

	//go:embed testdata/projects/dynamic_files/*
	// Validates That:
	//  1. Files are created
	//  2. Files are rendered
	//  3. Files are skipped
	//  4. Empty files are ignored
	dynamicFiles embed.FS

	//go:embed testdata/projects/injected_files/*
	injectedFiles embed.FS

	//go:embed testdata/projects/nested_scaffold/*
	// Validates That:
	//  1. Nested directories are created
	//  2. Nested files are created
	//  3. Nested files are rendered
	nestedFiles embed.FS

	//go:embed testdata/projects/invalid_project/*
	// Validates That:
	//  1. Invalid project structure is detected and error is returned
	invalidProject embed.FS

	//go:embed testdata/projects/feature_flag/*
	// Validates That:
	//  1. Conditional feature flags block files from being rendered
	featureFlag embed.FS
)

func FeatureFlagFiles() fs.FS {
	f, _ := fs.Sub(featureFlag, "testdata/projects/feature_flag")
	return f
}

func FeatureFlagProject() *Project {
	return &Project{
		NameTemplate: "{{ .ProjectKebab }}",
		Name:         "NewProject",
		Conf: &ProjectScaffoldFile{
			Features: []Feature{
				{
					Value: "{{ .Scaffold.feature }}",
					Globs: []string{
						"**/feature/**/*",
					},
				},
			},
		},
	}
}

func NestedFiles() fs.FS {
	f, _ := fs.Sub(nestedFiles, "testdata/projects/nested_scaffold")
	return f
}

func NestedFilesProject() *Project {
	return &Project{
		NameTemplate: "{{ .Project }}",
		Name:         "NewProject",
		Conf: &ProjectScaffoldFile{
			Computed: map[string]string{
				"snake_project": "{{ snakecase .Project }}",
			},
		},
	}
}

func DynamicFiles() fs.FS {
	f, _ := fs.Sub(dynamicFiles, "testdata/projects/dynamic_files")
	return f
}

func DynamicFilesProject() *Project {
	return &Project{
		NameTemplate: "{{ .Project }}",
		Name:         "NewProject",
		Conf: &ProjectScaffoldFile{
			Skip: []string{
				"copy.txt",
			},
		},
	}
}

func InjectedFiles() fs.FS {
	f, _ := fs.Sub(injectedFiles, "testdata/projects/injected_files")
	return f
}

func InjectedFilesProject() *Project {
	return &Project{
		NameTemplate: "{{ .Project }}",
		Name:         "NewProject",
		Conf: &ProjectScaffoldFile{
			Computed: map[string]string{
				"Site": "{{ .Project }}/site.yaml",
			},
			Inject: []Injectable{
				{
					Name:     "test1",
					Path:     "{{ .Computed.Site }}",
					At:       "# start",
					Template: "injected: true",
				},
			},
		},
	}
}

func InvalidStructure() fs.FS {
	f, _ := fs.Sub(invalidProject, "testdata/projects/invalid_project")
	return f
}

func InvalidStructureProject() *Project {
	return &Project{
		NameTemplate: "{{ .Project }}",
		Name:         "NewProject",
		Conf:         &ProjectScaffoldFile{},
	}
}

func CustomDelimsFiles() fs.FS {
	f, _ := fs.Sub(customDelimsFiles, "testdata/projects/custom_delims")
	return f
}

func CustomDelimsProject() *Project {
	return &Project{
		NameTemplate: "{{ .Project }}",
		Name:         "NewProject",
		Conf: &ProjectScaffoldFile{
			Delimiters: []Delimiters{
				{
					Glob:  "**/*custom.txt",
					Left:  "[[",
					Right: "]]",
				},
			},
		},
	}
}

func PartialsFiles() fs.FS {
	f, _ := fs.Sub(partialsFiles, "testdata/projects/with_partials")
	return f
}

func PartialsProject() *Project {
	return &Project{
		NameTemplate: "{{ .Project }}",
		Name:         "NewProject",
		Conf: &ProjectScaffoldFile{
			Partials: "partials",
			Computed: map[string]string{
				"Greeting": "Hello, World!",
			},
		},
	}
}

func TestMain(m *testing.M) {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.With().Logger().Level(zerolog.DebugLevel)

	os.Exit(m.Run())
}
