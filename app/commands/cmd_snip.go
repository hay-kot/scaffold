package commands

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"charm.land/glamour/v2"
	"charm.land/huh/v2"
	"github.com/hay-kot/scaffold/app/argparse"
	"github.com/hay-kot/scaffold/app/scaffold"
	"github.com/hay-kot/scaffold/app/scaffold/snippets"
	"github.com/hay-kot/scaffold/internal/styles"
	"github.com/rs/zerolog/log"
	"github.com/sahilm/fuzzy"
)

// SnipStdout is the --output value that writes the rendered snippet to stdout.
const SnipStdout = "-"

type FlagsSnip struct {
	// Output is a file path, an existing directory to write the snippet into,
	// or SnipStdout. Empty means the current directory.
	Output    string
	Overwrite bool
	NoPrompt  bool
	Preset    string
	List      bool
}

func (ctrl *Controller) Snip(args []string, flags FlagsSnip) error {
	ctrl.ready()

	roots := ctrl.snippetRoots()
	found := snippets.Discover(roots)

	if flags.List {
		return ctrl.listSnippets(found, roots)
	}

	if len(found) == 0 {
		return fmt.Errorf("no snippets found in %s", strings.Join(roots, ", "))
	}

	var name string
	if len(args) > 0 {
		name, args = args[0], args[1:]
	}

	snip, err := ctrl.pickSnippet(found, name, flags.NoPrompt)
	if err != nil {
		return err
	}

	argvars, err := argparse.Parse(args)
	if err != nil {
		return err
	}

	f, err := os.Open(snip.Path)
	if err != nil {
		return err
	}
	defer f.Close() //nolint:errcheck

	parsed, err := snippets.Parse(f)
	if err != nil {
		return fmt.Errorf("snippet %s: %w", snip.Name, err)
	}

	conf := parsed.Conf
	for _, key := range parsed.IgnoredKeys() {
		log.Warn().Str("snippet", snip.Name).Str("key", key).Msg("front matter key does not apply to snippets, ignoring it")
	}

	info, err := f.Stat()
	if err != nil {
		return err
	}

	// When the snippet goes to stdout, everything else must stay off stdout so
	// the output can be piped or redirected.
	var messages io.Writer = os.Stdout
	if flags.Output == SnipStdout {
		messages = os.Stderr
	}

	ok, err := conf.Metadata.IsCompatible(log.Logger, ctrl.Version)
	if err != nil {
		return err
	}
	if !ok {
		return fmt.Errorf("snippet requires version %s or higher", conf.Metadata.MinimumVersion)
	}

	if !flags.NoPrompt && conf.Messages.Pre != "" {
		out, err := glamour.RenderWithEnvironmentConfig(conf.Messages.Pre)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintln(messages, out)
	}

	p := &scaffold.Project{
		NameTemplate: scaffold.TemplateDirName,
		Conf:         conf,
	}

	vars, err := ctrl.snippetVars(p, argvars, flags, messages)
	if err != nil {
		return err
	}

	vars, err = scaffold.BuildVars(ctrl.engine, p, vars)
	if err != nil {
		return err
	}

	rendered, err := parsed.Render(ctrl.engine, snip.Name, snip.Path, vars)
	if err != nil {
		return err
	}

	if flags.Output == SnipStdout {
		if _, err := io.WriteString(os.Stdout, rendered); err != nil {
			return err
		}
	} else {
		if err := writeSnippet(flags.Output, filepath.Base(snip.Path), []byte(rendered), snippetPerm(info.Mode()), flags.Overwrite); err != nil {
			return err
		}
	}

	if !flags.NoPrompt && conf.Messages.Post != "" {
		post, err := ctrl.engine.TmplString(conf.Messages.Post, vars)
		if err != nil {
			return err
		}

		out, err := glamour.RenderWithEnvironmentConfig(post)
		if err != nil {
			return err
		}

		_, _ = fmt.Fprintln(messages, out)
	}

	return nil
}

func (ctrl *Controller) snippetVars(p *scaffold.Project, argvars map[string]any, flags FlagsSnip, messages io.Writer) (map[string]any, error) {
	if !flags.NoPrompt {
		input := scaffold.MergeMaps(ctrl.rc.Defaults, argvars)
		vars, err := p.AskQuestions(input, ctrl.engine, styles.Theme(ctrl.rc.Settings.Theme), messages)
		if err != nil {
			return nil, err
		}

		// AskQuestions names template scaffolds after their "templates"
		// directory. A snippet has no such directory, so keep the Project
		// value the user gave, the same as with --no-prompt.
		setSnippetProject(p, vars, input)
		return vars, nil
	}

	var preset map[string]any
	if flags.Preset != "" {
		var ok bool
		preset, ok = p.Conf.Presets[flags.Preset]
		if !ok {
			return nil, fmt.Errorf("preset '%s' not found", flags.Preset)
		}
	}

	vars := scaffold.MergeMaps(ctrl.rc.Defaults, preset, argvars)
	setSnippetProject(p, vars, vars)

	return vars, nil
}

func setSnippetProject(p *scaffold.Project, vars, input map[string]any) {
	name, _ := input["Project"].(string)
	p.Name = name

	if name == "" {
		delete(vars, "Project")
		return
	}

	vars["Project"] = name
}

// snippetPerm keeps the executable bits of the snippet source so that script
// snippets stay runnable, without passing through group or world write.
func snippetPerm(src fs.FileMode) fs.FileMode {
	if src.Perm()&0o111 != 0 {
		return 0o755
	}

	return 0o644
}

// snippetRoots returns the configured snippet locations followed by the ones
// passed with --snippet-dir. Without configured locations it falls back to a
// snippets directory next to the scaffoldrc file.
func (ctrl *Controller) snippetRoots() []string {
	rcDir := filepath.Dir(ctrl.Flags.ScaffoldRCPath)

	configured := ctrl.rc.Snippets
	if len(configured) == 0 {
		configured = []string{filepath.Join(rcDir, "snippets")}
	}

	roots := make([]string, 0, len(configured)+len(ctrl.Flags.SnippetDirs))
	for _, root := range configured {
		root = snippets.ExpandPath(root)
		if !filepath.IsAbs(root) {
			root = filepath.Join(rcDir, root)
		}
		roots = append(roots, root)
	}

	for _, root := range ctrl.Flags.SnippetDirs {
		roots = append(roots, snippets.ExpandPath(root))
	}

	return roots
}

func (ctrl *Controller) listSnippets(found []snippets.Snippet, roots []string) error {
	if len(found) == 0 {
		ctrl.printer.List("No snippets found in", roots)
		return nil
	}

	ctrl.printer.LineBreak()
	ctrl.printer.List("Snippets", snippets.Names(found))
	ctrl.printer.LineBreak()
	return nil
}

func (ctrl *Controller) pickSnippet(found []snippets.Snippet, name string, noPrompt bool) (snippets.Snippet, error) {
	if name == "" {
		if noPrompt {
			return snippets.Snippet{}, errors.New("snippet name is required, see 'scaffold snip --list' for available snippets")
		}

		var err error
		name, err = snippetPickerPrompt(snippets.Names(found), ctrl.rc.Settings.Theme)
		if err != nil {
			return snippets.Snippet{}, err
		}
	}

	snip, ok, err := snippets.Find(found, name)
	if err != nil {
		return snippets.Snippet{}, err
	}
	if ok {
		return snip, nil
	}

	if !noPrompt {
		matches := fuzzy.Find(name, snippets.Names(found))
		if len(matches) > 0 && didYouMeanPrompt(name, matches[0].Str, "snippet") {
			snip, _, err := snippets.Find(found, matches[0].Str)
			return snip, err
		}
	}

	return snippets.Snippet{}, fmt.Errorf("snippet %q not found, see 'scaffold snip --list' for available snippets", name)
}

func snippetPickerPrompt(names []string, theme styles.HuhTheme) (string, error) {
	options := make([]huh.Option[string], len(names))
	for i, name := range names {
		options[i] = huh.NewOption(name, name)
	}

	var selected string
	err := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a snippet").
				Options(options...).
				Filtering(true).
				Value(&selected),
		),
	).WithTheme(styles.Theme(theme)).Run()

	return selected, err
}

// writeSnippet writes content to dest. An empty
// dest, a dest that ends in a path separator, or an existing directory
// receives a file called filename. An existing file is only replaced when
// overwrite is set.
func writeSnippet(dest, filename string, content []byte, perm fs.FileMode, overwrite bool) error {
	if dest == "" {
		dest = "."
	}

	if strings.HasSuffix(dest, "/") || strings.HasSuffix(dest, string(filepath.Separator)) {
		dest = filepath.Join(dest, filename)
	} else if info, err := os.Stat(dest); err == nil && info.IsDir() {
		dest = filepath.Join(dest, filename)
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return err
	}

	flag := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	if !overwrite {
		flag = os.O_WRONLY | os.O_CREATE | os.O_EXCL
	}

	f, err := os.OpenFile(dest, flag, perm)
	if err != nil {
		if errors.Is(err, fs.ErrExist) {
			return fmt.Errorf("%s already exists, use --overwrite to replace it", dest)
		}
		return err
	}

	// OpenFile applies perm only to a new file, and an overwritten file must
	// still end up with the snippet's mode.
	if overwrite {
		if err := f.Chmod(perm); err != nil {
			_ = f.Close()
			return err
		}
	}

	if _, err := f.Write(content); err != nil {
		_ = f.Close()
		return err
	}

	return f.Close()
}
