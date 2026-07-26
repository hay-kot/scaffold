package pkgs

import (
	"archive/zip"
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	getter "github.com/hashicorp/go-getter"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoGetterDownloader_HTTPArchiveSubdirectory(t *testing.T) {
	var archive bytes.Buffer
	writer := zip.NewWriter(&archive)

	for name, contents := range map[string]string{
		"template-root/scaffold.yaml":       "---",
		"template-root/templates/hello.txt": "hello",
		"unrelated.txt":                     "ignore me",
	} {
		file, err := writer.Create(name)
		require.NoError(t, err)
		_, err = file.Write([]byte(contents))
		require.NoError(t, err)
	}
	require.NoError(t, writer.Close())

	server := httptest.NewServer(http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
		response.Header().Set("Content-Type", "application/zip")
		response.WriteHeader(http.StatusOK)
		if request.Method != http.MethodHead {
			_, _ = response.Write(archive.Bytes())
		}
	}))
	t.Cleanup(server.Close)

	destination := filepath.Join(t.TempDir(), "download")
	source := server.URL + "/template.zip//template-root"

	err := (goGetterDownloader{}).Download(destination, source)

	require.NoError(t, err)
	assert.FileExists(t, filepath.Join(destination, "scaffold.yaml"))
	assert.FileExists(t, filepath.Join(destination, "templates", "hello.txt"))
	assert.NoFileExists(t, filepath.Join(destination, "unrelated.txt"))

	mode, err := os.Stat(filepath.Join(destination, "templates", "hello.txt"))
	require.NoError(t, err)
	assert.False(t, mode.IsDir())
}

func TestNewGetters_DisablesTerraformRedirects(t *testing.T) {
	getters := newGetters()

	httpGetter, ok := getters["http"].(*getter.HttpGetter)
	require.True(t, ok)
	assert.True(t, httpGetter.XTerraformGetDisabled)
	assert.Same(t, getters["http"], getters["https"])
}

func TestIsGetterSource(t *testing.T) {
	for _, source := range []string{
		"http::https://example.com/template",
		"s3::https://s3.amazonaws.com/bucket/template",
		"gcs::https://storage.googleapis.com/bucket/template",
		"https://example.com/template.zip",
		"https://example.com/template.zip//nested",
		"https://github.com/example/template.git//nested",
		"https://example.com/template?archive=zip",
	} {
		t.Run(source, func(t *testing.T) {
			assert.True(t, isGetterSource(source))
		})
	}

	for _, source := range []string{
		"https://github.com/example/template",
		"git@github.com:example/template",
		"/local/template",
		"/local//template",
		"./local/template.zip",
	} {
		t.Run(source, func(t *testing.T) {
			assert.False(t, isGetterSource(source))
		})
	}
}
