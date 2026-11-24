// Package snippets discovers and parses single file scaffolds. A snippet is a
// template file with optional YAML front matter that uses the scaffold.yaml
// schema.
package snippets

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v4"
	"github.com/hay-kot/scaffold/app/core/apperrors"
	"github.com/hay-kot/scaffold/app/core/engine"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/rs/zerolog/log"
	"gopkg.in/yaml.v3"
)

type Snippet struct {
	// Name is the slash separated path of the snippet relative to the root it
	// was found in.
	Name string
	// Path is the location of the snippet file on disk.
	Path string
}

// ExpandPath expands environment variables and a leading ~ in p.
func ExpandPath(p string) string {
	p = os.ExpandEnv(p)
	if p == "~" || strings.HasPrefix(p, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			p = filepath.Join(home, p[1:])
		}
	}

	return p
}

func isGlob(root string) bool {
	return strings.ContainsAny(root, "*?[{")
}

// Discover returns the snippets found in roots. A root is a directory, which
// is searched recursively while following symlinks and skipping hidden files
// and directories, a single file, or a doublestar glob pattern. A path that
// exists is never read as a glob. Roots that do not exist are skipped, and
// roots that cannot be read are logged and skipped so one bad location does
// not hide the others. When two roots provide the same name, the earlier root
// wins.
func Discover(roots []string) []Snippet {
	seen := map[string]bool{}
	out := []Snippet{}

	for _, root := range roots {
		found, err := discoverRoot(root)
		if err != nil {
			log.Warn().Err(err).Str("root", root).Msg("skipping snippet root")
			continue
		}

		for _, s := range found {
			if seen[s.Name] {
				log.Debug().Str("name", s.Name).Str("path", s.Path).Msg("snippet shadowed by an earlier root, skipping")
				continue
			}

			seen[s.Name] = true
			out = append(out, s)
		}
	}

	return out
}

func discoverRoot(root string) ([]Snippet, error) {
	info, err := os.Stat(root)
	switch {
	case errors.Is(err, fs.ErrNotExist) && isGlob(root):
		return discoverGlob(root)
	case errors.Is(err, fs.ErrNotExist):
		log.Debug().Str("root", root).Msg("snippet root not found, skipping")
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed to check snippet root %s: %w", root, err)
	case !info.IsDir():
		return []Snippet{{Name: filepath.Base(root), Path: root}}, nil
	}

	return discoverDir(root)
}

// discoverDir walks root and follows symlinks, because dotfiles managers
// often link snippet directories into place. Each real directory is read
// once to stop symlink loops.
func discoverDir(root string) ([]Snippet, error) {
	out := []Snippet{}
	visited := map[string]bool{}

	var walk func(dir, prefix string) error
	walk = func(dir, prefix string) error {
		real, err := filepath.EvalSymlinks(dir)
		if err != nil {
			return err
		}
		if visited[real] {
			return nil
		}
		visited[real] = true

		entries, err := os.ReadDir(dir)
		if err != nil {
			return fmt.Errorf("failed to list snippets in %s: %w", dir, err)
		}

		for _, entry := range entries {
			if strings.HasPrefix(entry.Name(), ".") {
				continue
			}

			p := filepath.Join(dir, entry.Name())
			name := path.Join(prefix, entry.Name())

			info, err := os.Stat(p)
			if err != nil {
				log.Warn().Err(err).Str("path", p).Msg("skipping unreadable snippet")
				continue
			}

			switch {
			case info.IsDir():
				if err := walk(p, name); err != nil {
					return err
				}
			case info.Mode().IsRegular():
				out = append(out, Snippet{Name: name, Path: p})
			}
		}

		return nil
	}

	return out, walk(root, "")
}

func discoverGlob(root string) ([]Snippet, error) {
	base, pattern := doublestar.SplitPattern(filepath.ToSlash(root))

	if _, err := os.Stat(base); errors.Is(err, fs.ErrNotExist) {
		log.Debug().Str("root", root).Msg("snippet glob base not found, skipping")
		return nil, nil
	}

	matches, err := doublestar.Glob(os.DirFS(base), pattern, doublestar.WithFilesOnly())
	if err != nil {
		return nil, fmt.Errorf("failed to match snippet glob %s: %w", root, err)
	}

	sort.Strings(matches)

	out := make([]Snippet, 0, len(matches))
	for _, m := range matches {
		out = append(out, Snippet{Name: m, Path: filepath.Join(filepath.FromSlash(base), filepath.FromSlash(m))})
	}

	return out, nil
}

// Find returns the snippet called name. An exact match wins over a match that
// leaves out the file extension.
func Find(snippets []Snippet, name string) (Snippet, bool) {
	for _, s := range snippets {
		if s.Name == name {
			return s, true
		}
	}

	for _, s := range snippets {
		if strings.TrimSuffix(s.Name, path.Ext(s.Name)) == name {
			return s, true
		}
	}

	return Snippet{}, false
}

// Names returns the name of each snippet.
func Names(snippets []Snippet) []string {
	out := make([]string, len(snippets))
	for i, s := range snippets {
		out[i] = s.Name
	}

	return out
}

const frontMatterDelim = "---"

// File is a parsed snippet.
type File struct {
	Conf *scaffold.ProjectScaffoldFile
	Body string
	// BodyLine is the 1-based line of the snippet file where Body starts. It
	// maps template error lines in Body back to lines in the file.
	BodyLine int

	raw string
}

// Parse splits a snippet into its scaffold configuration and template body.
// The front matter is optional; it must start on the first line and end with
// a line that holds only "---". One blank line directly after the closing
// delimiter is dropped so authors can separate the two parts.
func Parse(r io.Reader) (*File, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	out := &File{Conf: &scaffold.ProjectScaffoldFile{}, BodyLine: 1, raw: string(data)}

	first, rest, ok := cutLine(data)
	if !ok || string(first) != frontMatterDelim {
		out.Body = string(data)
		return out, nil
	}

	consumed := 1
	var frontMatter []byte
	for {
		var line []byte
		line, rest, ok = cutLine(rest)
		consumed++
		if string(line) == frontMatterDelim {
			break
		}

		if !ok {
			return nil, errors.New("snippet front matter has no closing ---")
		}

		frontMatter = append(frontMatter, line...)
		frontMatter = append(frontMatter, '\n')
	}

	if err := yaml.Unmarshal(frontMatter, out.Conf); err != nil {
		return nil, fmt.Errorf("failed to parse snippet front matter: %w", err)
	}

	if blank, after, ok := cutLine(rest); ok && len(blank) == 0 {
		rest = after
		consumed++
	}

	out.Body = string(rest)
	out.BodyLine = consumed + 1
	return out, nil
}

// IgnoredKeys returns the front matter keys that are set but have no effect
// on a snippet, because they only make sense for a directory of files.
func (f *File) IgnoredKeys() []string {
	c := f.Conf
	keys := []string{}
	for _, k := range []struct {
		name string
		set  bool
	}{
		{"skip", len(c.Skip) > 0},
		{"rewrites", len(c.Rewrites) > 0},
		{"inject", len(c.Inject) > 0},
		{"features", len(c.Features) > 0},
		{"each", len(c.Each) > 0},
	} {
		if k.set {
			keys = append(keys, k.name)
		}
	}

	return keys
}

// Render executes the body with vars. The delimiters come from the last
// front matter delimiters entry whose glob matches name, the same rule a
// full scaffold applies to each file. Template errors report the snippet path
// and the line in the snippet file, with the lines around it.
func (f *File) Render(eng *engine.Engine, name, path string, vars any) (string, error) {
	left, right := "{{", "}}"
	for _, d := range f.Conf.Delimiters {
		match, err := doublestar.Match(d.Glob, name)
		if err != nil {
			return "", fmt.Errorf("invalid delimiters glob %q: %w", d.Glob, err)
		}

		if match {
			left, right = d.Left, d.Right
		}
	}

	tmpl, err := eng.Factory(strings.NewReader(f.Body), engine.WithDelims(left, right))
	if err != nil {
		if errors.Is(err, engine.ErrTemplateIsEmpty) {
			return "", nil
		}
		return "", f.templateError(err, path, left, right)
	}

	var out strings.Builder
	if err := tmpl.Execute(&out, vars); err != nil {
		return "", f.templateError(err, path, left, right)
	}

	return out.String(), nil
}

func (f *File) templateError(err error, path, left, right string) error {
	terr := apperrors.WrapTemplateError(err, path).WithDelimiters(left, right)
	if terr.LineNumber <= 0 {
		return terr
	}

	terr.LineNumber += f.BodyLine - 1

	lines := strings.Split(f.raw, "\n")
	if terr.LineNumber > len(lines) {
		return terr
	}

	from := max(terr.LineNumber-2, 0)
	to := min(terr.LineNumber+1, len(lines))
	return terr.WithContext(lines[from:to])
}

// cutLine returns the first line of data without its line ending, the data
// after it, and whether a line ending was found.
func cutLine(data []byte) (line, rest []byte, found bool) {
	line, rest, found = bytes.Cut(data, []byte("\n"))
	return bytes.TrimSuffix(line, []byte("\r")), rest, found
}
