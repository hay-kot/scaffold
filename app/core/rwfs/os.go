package rwfs

import (
	"context"
	"io/fs"
	"os"
	"os/exec"
	"os/signal"
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

func (o *OsWFS) RunHook(name string, data []byte, args []string) error {
	tmp, err := writeHook(name, data)

	defer func() {
		if rerr := os.Remove(tmp); rerr != nil && err == nil {
			err = rerr
		}
	}()

	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		// stop receiving signal notifications as soon as possible.
		<-ctx.Done()
		stop()
	}()

	cmd := exec.CommandContext(ctx, tmp, append([]string{tmp}, args...)...)
	cmd.Dir = o.root
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	return err
}

func writeHook(name string, data []byte) (string, error) {
	f, err := os.CreateTemp("", name)
	if err != nil {
		return "", err
	}

	tmp := f.Name()

	err = os.Chmod(tmp, 0700)
	if err != nil {
		return tmp, err
	}

	_, err = f.Write(data)
	if cerr := f.Close(); cerr != nil && err == nil {
		err = cerr
	}
	return tmp, err
}
