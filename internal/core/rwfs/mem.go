package rwfs

import (
	"io/fs"

	"github.com/psanford/memfs"
)

var _ WriteFS = &MemoryWFS{}
var _ fs.FS = &MemoryWFS{}

// MemoryWFS is a WFS implementation that uses a memory file system
// from github.com/psanford/memfs
// This should only be used for testing purposes, but may have other uses
// in the future.
type MemoryWFS struct {
	*memfs.FS
}

func NewMemoryWFS() *MemoryWFS {
	return &MemoryWFS{
		FS: memfs.New(),
	}
}
