package scaffold

import (
	"bytes"
	"io"
	"io/fs"
	"os"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/scaffold/internal/core/rwfs"
	"github.com/hay-kot/scaffold/internal/engine"
	"github.com/huandu/xstrings"
)

type RWFSArgs struct {
	ReadFS  rwfs.ReadFS
	WriteFS rwfs.WriteFS
	Project *Project
}

// RenderRWFS renders a rwfs.RFS to a rwfs.WriteFS by compiling all files in the rwfs.ReadFS
// and writing the compiled files to the WriteFS.
func RenderRWFS(s *engine.Engine, args *RWFSArgs, vars engine.Vars) error {
	iVars := engine.Vars{
		"Project":      args.Project.Name,
		"ProjectSnake": xstrings.ToSnakeCase(args.Project.Name),
		"ProjectKebab": xstrings.ToKebabCase(args.Project.Name),
		"ProjectCamel": xstrings.ToCamelCase(args.Project.Name),
		"Scaffold":     vars,
	}

	return fs.WalkDir(args.ReadFS, args.Project.NameTemplate, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if args.Project.Conf != nil && len(args.Project.Conf.Skip) > 0 {
			for _, pattern := range args.Project.Conf.Skip {
				// Use relative path for matching so that the config writers don't have to
				// specify every file as **/*.goreleaser.yml instead of just *.goreleaser.yml
				// to match things at the root of the project file.
				relativePath := strings.TrimPrefix(path, args.Project.NameTemplate+"/")
				match, err := doublestar.PathMatch(pattern, relativePath)
				if err != nil {
					return err
				}

				if !match {
					continue
				}

				outpath, err := s.TmplString(path, iVars)
				if err != nil {
					return err
				}

				rf, err := args.ReadFS.Open(path)

				if err != nil {
					return err
				}

				bits, err := io.ReadAll(rf)
				if err != nil {
					return err
				}

				err = args.WriteFS.WriteFile(outpath, bits, os.ModePerm)
				if err != nil {
					return err
				}

				return nil
			}
		}

		outpath, err := s.TmplString(path, iVars)
		if err != nil {
			return err
		}

		if d.IsDir() {
			err = args.WriteFS.MkdirAll(outpath, os.ModePerm)
			if err != nil {
				return err
			}

			return nil
		}

		f, err := args.ReadFS.Open(path)
		if err != nil {
			return err
		}

		tmpl, err := s.TmplFactory(f)
		if err != nil {
			_ = f.Close()
			return err
		}

		buff := bytes.NewBuffer(nil)

		err = s.RenderTemplate(buff, tmpl, iVars)
		if err != nil {
			_ = f.Close()
			return err
		}

		err = args.WriteFS.WriteFile(outpath, buff.Bytes(), os.ModePerm)
		if err != nil {
			_ = f.Close()
			return err
		}

		return f.Close()
	})
}
