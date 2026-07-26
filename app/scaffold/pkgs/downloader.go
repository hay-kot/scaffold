package pkgs

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"path/filepath"
	"strings"

	getter "github.com/hashicorp/go-getter"
)

type Downloader interface {
	Download(destination, source string) error
}

type DownloaderFunc func(destination, source string) error

func (f DownloaderFunc) Download(destination, source string) error {
	return f(destination, source)
}

type goGetterDownloader struct{}

func (goGetterDownloader) Download(destination, source string) error {
	client := &getter.Client{
		Ctx:             context.Background(),
		Src:             source,
		Dst:             destination,
		Mode:            getter.ClientModeDir,
		DisableSymlinks: true,
		Getters:         newGetters(),
	}

	return client.Get()
}

func newGetters() map[string]getter.Getter {
	httpGetter := &getter.HttpGetter{
		Netrc:                 true,
		XTerraformGetDisabled: true,
	}

	return map[string]getter.Getter{
		"file":  new(getter.FileGetter),
		"git":   new(getter.GitGetter),
		"gcs":   new(getter.GCSGetter),
		"hg":    new(getter.HgGetter),
		"s3":    new(getter.S3Getter),
		"http":  httpGetter,
		"https": httpGetter,
	}
}

func getterCacheDir(cache, source string) string {
	sum := sha256.Sum256([]byte(source))
	key := hex.EncodeToString(sum[:])

	return filepath.Join(cache, "_getter", key)
}

func isGetterSource(source string) bool {
	if separator := strings.Index(source, "::"); separator > 0 {
		protocol := source[:separator]
		if !strings.ContainsAny(protocol, `/\`) {
			return true
		}
	}

	root, subdir := getter.SourceDirSubdir(source)
	parsed, err := url.Parse(root)
	if err != nil {
		return false
	}

	switch scheme := strings.ToLower(parsed.Scheme); scheme {
	case "s3", "gcs":
		return true
	case "file", "git", "hg", "http", "https":
		// Supported by go-getter.
	default:
		return false
	}

	if subdir != "" {
		return true
	}

	if parsed.Query().Has("archive") {
		return true
	}

	path := strings.ToLower(parsed.Path)
	for _, extension := range []string{
		".zip",
		".tar",
		".tar.gz",
		".tgz",
		".tar.bz2",
		".tbz2",
		".tar.xz",
		".txz",
		".gz",
		".bz2",
		".xz",
	} {
		if strings.HasSuffix(path, extension) {
			return true
		}
	}

	return false
}
