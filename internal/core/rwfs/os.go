package rwfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

var _ WriteFS = &OsWFS{}

type OsWFS struct {
	fs.FS
	root string
}

// NewOsWFS returns a new OsWFS with the given root path
// The root path is used to join the path to the file system for
// all receiver operations
func NewOsWFS(root string) *OsWFS {
	return &OsWFS{
		FS:   os.DirFS(root),
		root: root,
	}
}

// MkdirAll wraps os.MkdirAll implementing and Joins the root path to the path
// before calling os.MkdirAll
func (o *OsWFS) MkdirAll(path string, perm fs.FileMode) error {
	return os.MkdirAll(filepath.Join(o.root, path), perm)
}

// WriteFile wraps os.WriteFile implementing and Joins the root path to the name/path
// before calling os.WriteFile
func (o *OsWFS) WriteFile(name string, data []byte, perm fs.FileMode) error {
	return os.WriteFile(filepath.Join(o.root, name), data, perm)
}
