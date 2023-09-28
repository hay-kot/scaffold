package pkgs

import (
	"io/fs"
	"path/filepath"
)

func ListLocal(f fs.FS) ([]string, error) {
	// .scaffold
	// └── model
	//     └── scaffold.yaml
	// └── controller
	//     └── scaffold.yaml

	outpaths := []string{}

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == "scaffold.yaml" || d.Name() == "scaffold.yml" {
			outpaths = append(outpaths, filepath.Base(filepath.Dir(path)))
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return outpaths, nil
}

// ListSystem traverses the filesystem and returns a list of all the package paths
// and references. This lists only the system scaffolds, and not the ones in the local
// .scaffold directory.
func ListSystem(f fs.FS) ([]string, error) {
	// Example Structure
	// Root
	// └── github.com
	//     └── hay-kot
	//	        └── scaffold-go-cli
	//			    └── repository files

	outpaths := []string{}

	// walk the file system for each directory in the root FS
	// stop when the directory contains a scaffold.yaml or scaffold.yml file
	// and add the path to the outpaths slice
	// Maximum recursion depth is
	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// check if maximum recursion depth is reached
		if d.IsDir() && filepath.Clean(path) != "." {
			depth := len(filepath.SplitList(path))
			if depth > 4 {
				return filepath.SkipDir
			}
		}

		// check if scaffold.yaml or scaffold.yml exists in the directory
		if d.Name() == "scaffold.yaml" || d.Name() == "scaffold.yml" {
			outpaths = append(outpaths, filepath.Dir(path))
			return filepath.SkipDir
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return outpaths, nil
}
