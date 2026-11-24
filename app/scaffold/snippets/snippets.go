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

// IsGlob reports whether root is a glob pattern rather than a plain path.
func IsGlob(root string) bool {
	return strings.ContainsAny(root, "*?[{")
}

// Discover returns the snippets found in roots. A root is a directory, which
// is searched recursively while skipping hidden files and directories, a
// single file, or a doublestar glob pattern. Roots that do not exist are
// skipped. When two roots provide the same name, the earlier root wins.
func Discover(roots []string) ([]Snippet, error) {
	seen := map[string]bool{}
	out := []Snippet{}

	for _, root := range roots {
		found, err := discoverRoot(root)
		if err != nil {
			return nil, err
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

	return out, nil
}

func discoverRoot(root string) ([]Snippet, error) {
	if IsGlob(root) {
		return discoverGlob(root)
	}

	info, err := os.Stat(root)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		log.Debug().Str("root", root).Msg("snippet root not found, skipping")
		return nil, nil
	case err != nil:
		return nil, fmt.Errorf("failed to check snippet root %s: %w", root, err)
	case !info.IsDir():
		return []Snippet{{Name: filepath.Base(root), Path: root}}, nil
	}

	out := []Snippet{}
	err = fs.WalkDir(os.DirFS(root), ".", func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if p != "." && strings.HasPrefix(d.Name(), ".") {
			if d.IsDir() {
				return fs.SkipDir
			}
			return nil
		}

		if d.IsDir() {
			return nil
		}

		out = append(out, Snippet{Name: p, Path: filepath.Join(root, filepath.FromSlash(p))})
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list snippets in %s: %w", root, err)
	}

	return out, nil
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

// Parse splits a snippet into its scaffold configuration and template body.
// The front matter is optional; it must start on the first line and end with
// a line that holds only "---". One blank line directly after the closing
// delimiter is dropped so authors can separate the two parts.
func Parse(r io.Reader) (*scaffold.ProjectScaffoldFile, string, error) {
	data, err := io.ReadAll(r)
	if err != nil {
		return nil, "", err
	}

	conf := &scaffold.ProjectScaffoldFile{}

	first, rest, ok := cutLine(data)
	if !ok || string(first) != frontMatterDelim {
		return conf, string(data), nil
	}

	var frontMatter []byte
	for {
		var line []byte
		line, rest, ok = cutLine(rest)
		if string(line) == frontMatterDelim {
			break
		}

		if !ok {
			return nil, "", errors.New("snippet front matter has no closing ---")
		}

		frontMatter = append(frontMatter, line...)
		frontMatter = append(frontMatter, '\n')
	}

	if err := yaml.Unmarshal(frontMatter, conf); err != nil {
		return nil, "", fmt.Errorf("failed to parse snippet front matter: %w", err)
	}

	if blank, after, ok := cutLine(rest); ok && len(blank) == 0 {
		rest = after
	}

	return conf, string(rest), nil
}

// cutLine returns the first line of data without its line ending, the data
// after it, and whether a line ending was found.
func cutLine(data []byte) (line, rest []byte, found bool) {
	line, rest, found = bytes.Cut(data, []byte("\n"))
	return bytes.TrimSuffix(line, []byte("\r")), rest, found
}
