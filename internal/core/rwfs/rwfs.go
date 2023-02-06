// Package rwfs provides a Read Write File System to extend the standard
// io/fs package with the ability to write to the file system.
package rwfs

import (
	"io"
	"io/fs"
)

// ReadFS is a read only file system that can be used to read files from
// a file system. It is a alias for fs.FS.
type ReadFS = fs.FS

// WFile is a file that can be written to. It is a alias for fs.File.
// that also implements io.Writer.
type WFile interface {
	fs.File
	io.Writer
}

// WriteFS is a file system that can be used to read and write files.
// It is a alias for fs.FS that also implements Mkdir, MkdirAll and Create.
type WriteFS interface {
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
}
