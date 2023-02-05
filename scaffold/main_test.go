package scaffold

import (
	"embed"
	"io/fs"
	"os"
	"testing"

	"github.com/rs/zerolog"
)

var (
	//go:embed testdata/dynamic_files/*
	dynamicFiles embed.FS

	//go:embed testdata/nested_scaffold/*
	nestedFiles embed.FS

	//go:embed testdata/invalid_project/*
	invalidProject embed.FS
)

func NestedFiles() fs.FS {
	f, _ := fs.Sub(nestedFiles, "testdata/nested_scaffold")
	return f
}

func DynamicFiles() fs.FS {
	f, _ := fs.Sub(dynamicFiles, "testdata/dynamic_files")
	return f
}

func InvalidProject() fs.FS {
	f, _ := fs.Sub(invalidProject, "testdata/invalid_project")
	return f
}

func TestMain(m *testing.M) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)

	os.Exit(m.Run())
}
