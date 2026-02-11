package scaffold

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/scaffold/app/core/apperrors"
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

// eachPattern matches [varname] in path segments. The varname must be a valid
// identifier (letters, digits, underscores).
var eachPattern = regexp.MustCompile(`\[([a-zA-Z_]\w*)\]`)

// detectEachPattern finds the first [varname] in a path and returns the variable
// name and the bracket-wrapped token (e.g., "[services]"). Returns empty strings
// if no pattern is found.
func detectEachPattern(path string) (varName, token string) {
	m := eachPattern.FindStringSubmatch(path)
	if m == nil {
		return "", ""
	}
	return m[1], m[0]
}

// isEachVar checks whether varName is declared in the each config.
func isEachVar(each []EachConfig, varName string) (EachConfig, bool) {
	for _, ec := range each {
		if ec.Var == varName {
			return ec, true
		}
	}
	return EachConfig{}, false
}

// resolveListVar extracts a []string from vars["Scaffold"][varName].
// If the value is a plain string, it's wrapped into a single-element slice.
func resolveListVar(vars engine.Vars, varName string) ([]string, error) {
	scaffold, ok := vars["Scaffold"]
	if !ok {
		return nil, fmt.Errorf("each: Scaffold vars not found")
	}

	scaffoldMap, ok := scaffold.(engine.Vars)
	if !ok {
		return nil, fmt.Errorf("each: Scaffold vars is not a map")
	}

	val, ok := scaffoldMap[varName]
	if !ok {
		return nil, fmt.Errorf("each: variable %q not found in Scaffold vars", varName)
	}

	switch v := val.(type) {
	case []string:
		return v, nil
	case []any:
		out := make([]string, len(v))
		for i, item := range v {
			s, ok := item.(string)
			if !ok {
				return nil, fmt.Errorf("each: variable %q item %d is not a string", varName, i)
			}
			out[i] = s
		}
		return out, nil
	case string:
		return []string{v}, nil
	default:
		return nil, fmt.Errorf("each: variable %q has unsupported type %T", varName, val)
	}
}

// addFileContextToError reads the file content and adds context lines to the template error
func addFileContextToError(terr *apperrors.TemplateError, rfs rwfs.ReadFS, path string) *apperrors.TemplateError {
	if terr.LineNumber <= 0 {
		return terr
	}

	content, err := fs.ReadFile(rfs, path)
	if err != nil {
		return terr
	}

	lines := strings.Split(string(content), "\n")
	if terr.LineNumber > len(lines) {
		return terr
	}

	contextLines := []string{}

	// Line before (if exists)
	if terr.LineNumber > 1 {
		contextLines = append(contextLines, lines[terr.LineNumber-2])
	}

	// Error line
	contextLines = append(contextLines, lines[terr.LineNumber-1])

	// Line after (if exists)
	if terr.LineNumber < len(lines) {
		contextLines = append(contextLines, lines[terr.LineNumber])
	}

	return terr.WithContext(contextLines)
}

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

type processFileArgs struct {
	sourcePath string
	outpath    string
	d          fs.DirEntry
	vars       engine.Vars
}

func processFile(eng *engine.Engine, args *RWFSArgs, pf processFileArgs) error {
	f, err := args.ReadFS.Open(pf.sourcePath)
	if err != nil {
		log.Debug().Err(err).Str("path", pf.sourcePath).Msg("failed to open file")
		return err
	}

	delimLeft := "{{"
	delimRight := "}}"

	for _, delimOverride := range args.Project.Conf.Delimiters {
		relativePath := strings.TrimPrefix(pf.sourcePath, args.Project.NameTemplate+"/")
		match, err := doublestar.Match(delimOverride.Glob, relativePath)
		if err != nil {
			_ = f.Close()
			return err
		}

		if !match {
			continue
		}

		log.Debug().Str("outputh", pf.outpath).Str("glob", delimOverride.Glob).Msg("matched delimiter override")

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

		terr := apperrors.WrapTemplateError(err, pf.sourcePath).WithDelimiters(delimLeft, delimRight)
		terr = addFileContextToError(terr, args.ReadFS, pf.sourcePath)
		return terr
	}

	buff := bytes.NewBuffer(nil)

	err = tmpl.Execute(buff, pf.vars)
	if err != nil {
		_ = f.Close()
		terr := apperrors.WrapTemplateError(err, pf.sourcePath).WithDelimiters(delimLeft, delimRight)
		terr = addFileContextToError(terr, args.ReadFS, pf.sourcePath)
		return terr
	}

	if buff.Len() == 0 {
		_ = f.Close()
		return nil
	}

	if len(strings.TrimSpace(buff.String())) == 0 {
		_ = f.Close()
		return nil
	}

	err = args.WriteFS.MkdirAll(filepath.Dir(pf.outpath), os.ModePerm)
	if err != nil {
		if !os.IsExist(err) {
			_ = f.Close()
			return err
		}
	}

	err = args.WriteFS.WriteFile(pf.outpath, buff.Bytes(), os.ModePerm)
	if err != nil {
		_ = f.Close()
		return err
	}

	return f.Close()
}

// expandEachDir walks a source directory tree and calls processFile for each
// non-directory entry, replacing [token] in paths with the replacement value.
func expandEachDir(eng *engine.Engine, args *RWFSArgs, guards []filepathGuard, sourceDirPath string, token string, replacement string, vars engine.Vars) error {
	return fs.WalkDir(args.ReadFS, sourceDirPath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		outpath := strings.Replace(path, token, replacement, 1)

		for i, guard := range guards {
			outpath, err = guard(outpath, d)
			if err != nil {
				if errors.Is(err, errSkipRender) || errors.Is(err, errSkipWrite) {
					return nil
				}
				log.Debug().Err(err).Str("outpath", outpath).Int("guard", i).Msg("guard failed")
				return err
			}
		}

		if args.Project.NameTemplate == TemplateDirName {
			outpath = strings.TrimPrefix(outpath, TemplateDirName+"/")
		}

		return processFile(eng, args, processFileArgs{
			sourcePath: path,
			outpath:    outpath,
			d:          d,
			vars:       vars,
		})
	})
}

// makeEachVars creates a copy of vars with .Each set for the current iteration item.
func makeEachVars(vars engine.Vars, item string, index int) engine.Vars {
	v := maps.Clone(vars)
	v["Each"] = map[string]any{
		"Item":  item,
		"Index": index,
	}
	return v
}

// RenderRWFS renders a rwfs.RFS to a rwfs.WriteFS by compiling all files in the rwfs.ReadFS
// and writing the compiled files to the WriteFS.
func RenderRWFS(eng *engine.Engine, args *RWFSArgs, vars engine.Vars) error {
	const PartialsDir = "partials"

	// Path guards apply to all files (skipped and rendered)
	rewriteGuard := guardRewrite(args)
	renderPathGuard := guardRenderPath(eng, vars)
	noClobberGuard := guardNoClobber(args)

	pathGuards := []filepathGuard{
		rewriteGuard,
		renderPathGuard,
		noClobberGuard,
	}

	// Full guard chain for rendered files only
	guards := []filepathGuard{
		rewriteGuard,
		renderPathGuard,
		noClobberGuard,
		guardDirectories(args),
		guardFeatureFlag(eng, args, vars),
	}

	_, err := args.ReadFS.Open(PartialsDir)
	if err == nil {
		partialsFS, err := fs.Sub(args.ReadFS, PartialsDir)
		if err != nil {
			return fmt.Errorf("failed to create partials FS: %w", err)
		}

		err = eng.RegisterPartialsFS(partialsFS, ".")
		if err != nil {
			return fmt.Errorf("failed to register partials FS: %w", err)
		}
	}

	// Build a set of each-declared variable names for fast lookup.
	eachConfigs := args.Project.Conf.Each

	// Track directories that are being expanded so we skip their children
	// during the normal walk (they're handled by expandEachDir).
	expandedDirs := map[string]bool{}

	err = fs.WalkDir(args.ReadFS, args.Project.NameTemplate, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// If this path is under an already-expanded directory, skip it.
		for dir := range expandedDirs {
			if strings.HasPrefix(path, dir+"/") || path == dir {
				return nil
			}
		}

		// Detect [var] pattern in the path
		varName, token := detectEachPattern(path)
		if varName != "" {
			if ec, ok := isEachVar(eachConfigs, varName); ok {
				items, err := resolveListVar(vars, varName)
				if err != nil {
					return err
				}

				if len(items) == 0 {
					log.Warn().Str("var", varName).Msg("each variable is empty, no files generated")
					if d.IsDir() {
						return fs.SkipDir
					}
					return nil
				}

				if d.IsDir() {
					expandedDirs[path] = true

					for i, item := range items {
						iterVars := makeEachVars(vars, item, i)

						replacement := item
						if ec.As != "" {
							replacement, err = eng.TmplString(ec.As, iterVars)
							if err != nil {
								return fmt.Errorf("each: rendering 'as' template for %q: %w", varName, err)
							}
							if replacement == "" {
								return fmt.Errorf("each: 'as' template for %q rendered empty string", varName)
							}
							if strings.Contains(replacement, "/") {
								return fmt.Errorf("each: 'as' template for %q rendered value containing path separator: %q", varName, replacement)
							}
						}

						// Need per-iteration guards that use iterVars for path rendering
						iterRenderPathGuard := guardRenderPath(eng, iterVars)
						iterGuards := []filepathGuard{
							rewriteGuard,
							iterRenderPathGuard,
							noClobberGuard,
							guardDirectories(args),
							guardFeatureFlag(eng, args, iterVars),
						}

						if err := expandEachDir(eng, args, iterGuards, path, token, replacement, iterVars); err != nil {
							return err
						}
					}

					return fs.SkipDir
				}

				// File-level expansion
				for i, item := range items {
					iterVars := makeEachVars(vars, item, i)

					replacement := item
					if ec.As != "" {
						replacement, err = eng.TmplString(ec.As, iterVars)
						if err != nil {
							return fmt.Errorf("each: rendering 'as' template for %q: %w", varName, err)
						}
						if replacement == "" {
							return fmt.Errorf("each: 'as' template for %q rendered empty string", varName)
						}
						if strings.Contains(replacement, "/") {
							return fmt.Errorf("each: 'as' template for %q rendered value containing path separator: %q", varName, replacement)
						}
					}

					outpath := strings.Replace(path, token, replacement, 1)

					iterRenderPathGuard := guardRenderPath(eng, iterVars)
					iterGuards := []filepathGuard{
						rewriteGuard,
						iterRenderPathGuard,
						noClobberGuard,
						guardDirectories(args),
						guardFeatureFlag(eng, args, iterVars),
					}

					for gi, guard := range iterGuards {
						outpath, err = guard(outpath, d)
						if err != nil {
							if errors.Is(err, errSkipRender) || errors.Is(err, errSkipWrite) {
								goto nextFileItem
							}
							log.Debug().Err(err).Str("outpath", outpath).Int("guard", gi).Msg("guard failed")
							return err
						}
					}

					if args.Project.NameTemplate == TemplateDirName {
						outpath = strings.TrimPrefix(outpath, TemplateDirName+"/")
					}

					if err := processFile(eng, args, processFileArgs{
						sourcePath: path,
						outpath:    outpath,
						d:          d,
						vars:       iterVars,
					}); err != nil {
						return err
					}

				nextFileItem:
				}

				return nil
			}
		}

		// --- Normal (non-expanded) path ---

		if args.Project.Conf != nil && len(args.Project.Conf.Skip) > 0 {
			relativePath := strings.TrimPrefix(path, args.Project.NameTemplate+"/")

			for _, pattern := range args.Project.Conf.Skip {
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

				outpath := path
				for _, guard := range pathGuards {
					outpath, err = guard(outpath, d)
					if err != nil {
						if errors.Is(err, errFileExists) {
							return nil
						}
						return err
					}
				}

				if args.Project.NameTemplate == TemplateDirName {
					outpath = strings.TrimPrefix(outpath, TemplateDirName+"/")
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

		if args.Project.NameTemplate == TemplateDirName {
			outpath = strings.TrimPrefix(outpath, TemplateDirName+"/")
		}

		return processFile(eng, args, processFileArgs{
			sourcePath: path,
			outpath:    outpath,
			d:          d,
			vars:       vars,
		})
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
