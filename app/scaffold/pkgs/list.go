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

	// Ensure that the .scaffold directory exists
	// If it doesn't, return an empty slice
	if _, err := fs.Stat(f, ".scaffold"); err != nil {
		return outpaths, nil
	}

	err := fs.WalkDir(f, ".scaffold", func(path string, d fs.DirEntry, err error) error {
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

type PackageList struct {
	Root        string
	SubPackages []string
}

// ListSystem traverses the filesystem and returns a list of all the package paths
// and references. This lists only the system scaffolds, and not the ones in the local
// .scaffold directory.
func ListSystem(f fs.FS) ([]PackageList, error) {
	// Example Structure
	// Root
	// └── github.com
	//     └── hay-kot
	//	        └── scaffold-go-cli
	//			    └── repository files

	pkgs := []PackageList{}

	// walk the file system for each directory in the root FS
	// stop when the directory contains a scaffold.yaml or scaffold.yml file
	// and add the path to the outpaths slice
	// Maximum recursion depth is
	// Root
	current := PackageList{}

	err := fs.WalkDir(f, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			// check if ".git" directory exists in the directory and set as root
			if d.Name() == ".git" {
				if current.Root != "" {
					pkgs = append(pkgs, current)
				}

				current = PackageList{
					Root: filepath.Dir(path),
				}
			}
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
			// if not in the root directory, add the directory to the subpackages
			if current.Root != filepath.Dir(path) {
				current.SubPackages = append(current.SubPackages, filepath.Base(filepath.Dir(path)))
				return filepath.SkipDir
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	if current.Root != "" {
		pkgs = append(pkgs, current)
	}

	return pkgs, nil
}
