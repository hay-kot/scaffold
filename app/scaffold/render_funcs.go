package scaffold

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/core/rwfs"
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
var (
	errSkipRender = errors.New("skip render")
	errFileExists = errors.New("file exists and no-clobber is set to true")
	errSkipWrite  = errors.New("skip write")
)

type filepathGuard func(outpath string, f fs.DirEntry) (newOutpath string, err error)

func guardNoOp(outpath string, f fs.DirEntry) (string, error) {
	return outpath, nil
}

func guardRewrite(args *RWFSArgs) filepathGuard {
	if len(args.Project.Conf.Rewrites) == 0 {
		return guardNoOp
	}

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
	if !args.Project.Options.NoClobber {
		return guardNoOp
	}

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

		// if a directory matches a root level directory in the project
		// it is created, otherwise it is skipped last entry is skipped
		// because it is the "templates" directory
		for _, s := range projectNames {
			if outpath == s {
				err := args.WriteFS.MkdirAll(outpath, os.ModePerm)
				if err != nil {
					if !os.IsExist(err) {
						log.Debug().Err(err).Str("path", outpath).Msg("failed to create directory")
						return "", err
					}
				}
			}
		}

		return "", errSkipWrite
	}
}

func guardFeatureFlag(e *engine.Engine, args *RWFSArgs, vars engine.Vars) filepathGuard {
	if len(args.Project.Conf.Features) == 0 {
		return guardNoOp
	}

	return func(outpath string, f fs.DirEntry) (newOutpath string, err error) {
		for _, feature := range args.Project.Conf.Features {
			render, err := e.TmplString(feature.Value, vars)
			if err != nil {
				return "", err
			}

			booly, _ := strconv.ParseBool(render)
			if !booly {
				for _, pattern := range feature.Globs {
					match, err := doublestar.Match(pattern, outpath)
					if err != nil {
						log.Debug().Err(err).Str("path", outpath).Str("pattern", pattern).Msg("feature pattern match")
						return "", err
					}

					if match {
						return "", errSkipRender
					}
				}
			}
		}

		return outpath, nil
	}
}

// BuildVars builds the vars for the engine by setting the provided vars
// under the "Scaffold" key and adding the project name and computed vars.
func BuildVars(eng *engine.Engine, project *Project, vars engine.Vars) (engine.Vars, error) {
	iVars := engine.Vars{
		"Project":       project.Name,
		"ProjectSnake":  xstrings.ToSnakeCase(project.Name),
		"ProjectKebab":  xstrings.ToKebabCase(project.Name),
		"ProjectCamel":  xstrings.ToCamelCase(project.Name),
		"ProjectPascal": xstrings.ToPascalCase(project.Name),
		"Scaffold":      vars,
	}

	computed := make(map[string]any, len(project.Conf.Computed))
	for k, v := range project.Conf.Computed {
		out, err := eng.TmplString(v, iVars)
		if err != nil {
			return nil, err
		}

		// We must parse the integer first to avoid incorrectly parsing a '0'
		// as a boolean.
		if i, err := strconv.Atoi(out); err == nil {
			computed[k] = i
			continue
		}

		if b, err := strconv.ParseBool(out); err == nil {
			computed[k] = b
			continue
		}

		computed[k] = out
	}

	iVars["Computed"] = computed

	return iVars, nil
}

// RenderRWFS renders a rwfs.RFS to a rwfs.WriteFS by compiling all files in the rwfs.ReadFS
// and writing the compiled files to the WriteFS.
func RenderRWFS(eng *engine.Engine, args *RWFSArgs, vars engine.Vars) error {
	guards := []filepathGuard{
		guardRewrite(args),
		guardRenderPath(eng, vars),
		guardNoClobber(args),
		guardDirectories(args),
		guardFeatureFlag(eng, args, vars),
	}

	err := fs.WalkDir(args.ReadFS, args.Project.NameTemplate, func(path string, d fs.DirEntry, err error) error {
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

				rf, err := args.ReadFS.Open(path)
				if err != nil {
					return err
				}

				outpath, err := eng.TmplString(path, vars)
				if err != nil {
					return err
				}

				err = args.WriteFS.MkdirAll(filepath.Dir(outpath), os.ModePerm)
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

		delimLeft := "{{"
		delimRight := "}}"

		for _, delimOverride := range args.Project.Conf.Delimiters {
			// Use relative path for matching so that the config writers don't have to
			// specify every file as **/*.goreleaser.yml instead of just *.goreleaser.yml
			// to match things at the root of the project file.
			relativePath := strings.TrimPrefix(path, args.Project.NameTemplate+"/")
			match, err := doublestar.Match(delimOverride.Glob, relativePath)
			if err != nil {
				_ = f.Close()
				return err
			}

			if !match {
				continue
			}

			log.Debug().Str("outputh", outpath).Str("glob", delimOverride.Glob).Msg("matched delimiter override")

			if delimOverride.Left == "" || delimOverride.Right == "" {
				log.Error().
					Str("left", delimOverride.Left).
					Str("right", delimOverride.Right).
					Msg("override delimiters must not be empty")
			}

			delimLeft = delimOverride.Left
			delimRight = delimOverride.Right
		}

		tmpl, err := eng.Factory(f, engine.WithDelims(delimLeft, delimRight))
		if err != nil {
			_ = f.Close()

			if errors.Is(err, engine.ErrTemplateIsEmpty) {
				return nil
			}

			return err
		}

		buff := bytes.NewBuffer(nil)

		err = tmpl.Execute(buff, vars)
		if err != nil {
			_ = f.Close()
			return err
		}

		// Skip empty files
		if buff.Len() == 0 {
			_ = f.Close()
			return nil
		}

		// Skip whitespace files
		if len(strings.TrimSpace(buff.String())) == 0 {
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
	if err != nil {
		return err
	}

	// Do Injection Jobs

	for _, injection := range args.Project.Conf.Inject {
		path, err := eng.TmplString(injection.Path, vars)
		if err != nil {
			return err
		}

		f, err := args.WriteFS.Open(path)
		if err != nil {
			return err
		}

		out, err := eng.TmplString(injection.Template, vars)
		if err != nil {
			return err
		}

		// Assume that empty string or only whitespace is not a valid injection.
		if out == "" || strings.TrimSpace(out) == "" {
			continue
		}

		outbytes, err := Inject(f, out, injection.At, injection.Mode)
		if err != nil {
			return err
		}

		err = args.WriteFS.WriteFile(path, outbytes, os.ModePerm)
		if err != nil {
			return err
		}
	}

	return nil
}
