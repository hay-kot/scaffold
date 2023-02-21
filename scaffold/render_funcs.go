package scaffold

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/scaffold/internal/core/rwfs"
	"github.com/hay-kot/scaffold/internal/engine"
	"github.com/huandu/xstrings"
	"github.com/rs/zerolog/log"
)

type RWFSArgs struct {
	ReadFS  rwfs.ReadFS
	WriteFS rwfs.WriteFS
	Project *Project
}

// errSkipRender is used to skip rendering a file when a guard returns it.
// this should only be used in guards.
var errSkipRender = errors.New("skip render")
var errFileExists = errors.New("file exists and no-clobber is set to true")
var errSkipWrite = errors.New("skip write")

type filepathGuard func(outpath string, f fs.DirEntry) (newOutpath string, err error)

func guardRewrite(args *RWFSArgs) filepathGuard {
	return func(path string, f fs.DirEntry) (string, error) {
		for _, rewrite := range args.Project.Conf.Rewrites {
			match, err := doublestar.Match(rewrite.From, path)
			if err != nil {
				log.Debug().Err(err).Str("path", path).Str("pattern", rewrite.From).Msg("rewrite pattern match")
				return "", err
			}

			if match {
				return rewrite.To, nil
			}
		}

		return path, nil
	}
}

func guardRenderPath(s *engine.Engine, vars any) filepathGuard {
	return func(outpath string, f fs.DirEntry) (string, error) {
		outpath, err := s.TmplString(outpath, vars)
		if err != nil {
			log.Debug().Err(err).Str("path", outpath).Msg("failed to render project path")
			return "", err
		}

		return outpath, nil
	}
}

func guardNoClobber(args *RWFSArgs) filepathGuard {
	return func(outpath string, f fs.DirEntry) (string, error) {
		wf, err := args.WriteFS.Open(outpath)

		if err == nil {
			_ = wf.Close()
			if args.Project.Options.NoClobber {
				log.Debug().Str("path", outpath).Msg("file exists and no-clobber is set to true")
				return "", errFileExists
			}
		}

		return outpath, nil
	}
}

func guardDirectories(args *RWFSArgs) filepathGuard {
	return func(outpath string, f fs.DirEntry) (string, error) {
		if !f.IsDir() {
			return outpath, nil
		}

		// skip "/templates" directory
		match, _ := filepath.Match("templates", outpath)
		if match {
			return "", errSkipRender
		}

		err := args.WriteFS.MkdirAll(outpath, os.ModePerm)
		if err != nil {
			if !os.IsExist(err) {
				log.Debug().Err(err).Str("path", outpath).Msg("failed to create directory")
				return "", err
			}
		}

		return "", errSkipWrite
	}
}

// RenderRWFS renders a rwfs.RFS to a rwfs.WriteFS by compiling all files in the rwfs.ReadFS
// and writing the compiled files to the WriteFS.
func RenderRWFS(eng *engine.Engine, args *RWFSArgs, vars engine.Vars) error {
	iVars := engine.Vars{
		"Project":      args.Project.Name,
		"ProjectSnake": xstrings.ToSnakeCase(args.Project.Name),
		"ProjectKebab": xstrings.ToKebabCase(args.Project.Name),
		"ProjectCamel": xstrings.ToCamelCase(args.Project.Name),
		"Scaffold":     vars,
	}

	guards := []filepathGuard{
		guardRewrite(args),
		guardRenderPath(eng, iVars),
		guardNoClobber(args),
		guardDirectories(args),
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

				log.Debug().Str("path", path).Str("pattern", pattern).Bool("match", match).Msg("skip pattern match")
				if !match {
					continue
				}

				outpath, err := eng.TmplString(path, iVars)
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

		outpath := path

		for i, guard := range guards {
			outpath, err = guard(outpath, d)
			if err != nil {
				if errors.Is(err, errSkipRender) || errors.Is(err, errSkipWrite) {
					return nil
				}

				log.Debug().Err(err).Str("outpath", outpath).Int("guard", i).Msg("guard failed")
				return err
			}

			log.Debug().Str("outpath", outpath).Int("guard", i).Msg("guard")
		}

		f, err := args.ReadFS.Open(path)
		if err != nil {
			log.Debug().Err(err).Str("path", path).Msg("failed to open file")
			return err
		}

		tmpl, err := eng.Factory(f)
		if err != nil {
			_ = f.Close()

			if errors.Is(err, engine.ErrTemplateIsEmpty) {
				return nil
			}

			return err
		}

		buff := bytes.NewBuffer(nil)

		err = eng.Render(buff, tmpl, iVars)
		if err != nil {
			_ = f.Close()
			return err
		}

		// Skip empty files
		if buff.Len() == 0 {
			_ = f.Close()
			return nil
		}

		// Ensure the directory exists
		err = args.WriteFS.MkdirAll(filepath.Dir(outpath), os.ModePerm)
		if err != nil {
			if !os.IsExist(err) {
				_ = f.Close()
				return err
			}
		}

		err = args.WriteFS.WriteFile(outpath, buff.Bytes(), os.ModePerm)
		if err != nil {
			_ = f.Close()
			return err
		}

		return f.Close()
	})
}
