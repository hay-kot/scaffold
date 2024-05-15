// Package rwfs provides a Read Write File System to extend the standard
// io/fs package with the ability to write to the file system.
package rwfs

import (
	"errors"
	"io/fs"
)

var ErrHooksNotSupported = errors.New("hooks not supported")

// ReadFS is a read only file system that can be used to read files from
// a file system. It is a alias for fs.FS.
type ReadFS = fs.FS

// WriteFS is a file system that can be used to read and write files.
// It is a alias for fs.FS that also implements Mkdir, MkdirAll and Create.
type WriteFS interface {
	fs.FS
	MkdirAll(path string, perm fs.FileMode) error
	WriteFile(name string, data []byte, perm fs.FileMode) error
	RunHook(name string, data []byte, args []string) error
}
