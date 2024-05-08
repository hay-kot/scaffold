package rwfs

import (
	"io/fs"
	"strings"

	"github.com/psanford/memfs"
)

var (
	_ WriteFS = &MemoryWFS{}
	_ fs.FS   = &MemoryWFS{}
)

// MemoryWFS is a WFS implementation that uses a memory file system
// from github.com/psanford/memfs
// This should only be used for testing purposes, but may have other uses
// in the future.
type MemoryWFS struct {
	*memfs.FS
}

func (m *MemoryWFS) MkdirAll(path string, perm fs.FileMode) error {
	if path == "/" {
		// special case root dir always exists
		return nil
	}

	return m.FS.MkdirAll(path, perm)
}

func (m *MemoryWFS) WriteFile(path string, data []byte, perm fs.FileMode) error {
	path = strings.TrimPrefix(path, "/")
	return m.FS.WriteFile(path, data, perm)
}

func NewMemoryWFS() *MemoryWFS {
	return &MemoryWFS{
		FS: memfs.New(),
	}
}
