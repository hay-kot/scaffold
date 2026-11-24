---
---

# Snippets

A snippet is a single file template. Use one when a full scaffold with a `scaffold.yaml` and a template directory is more than you need, for example a `Taskfile.yml`, a `.gitignore`, or a README.

## File Format

A snippet is the template body, optionally preceded by YAML front matter. The front matter uses the same schema as [`scaffold.yaml`](/configuration/scaffold-file), and it starts and ends with a line that holds only `---`.

```markdown
---
messages:
  pre: |
    # README snippet
  post: |
    Wrote a README for {{ .Scaffold.name }}.

questions:
  - name: "name"
    prompt:
      message: "Project name"
    required: true
  - name: "description"
    prompt:
      message: "Description of the project"
---

# {{ .Scaffold.name }}

{{ .Scaffold.description }}
```

One blank line directly after the closing `---` is dropped, so you can separate the front matter from the body. A file without front matter is rendered as is.

These front matter keys apply to snippets:

- `questions`
- `computed`
- `messages`
- `presets`
- `metadata`
- `delimiters`, matched against the snippet name (for example `workflows/*.yml`)

Hooks, partials, `skip`, `rewrites`, `inject`, `features`, and `each` apply only to full scaffolds. Scaffold prints a warning when a snippet sets one of these keys, and ignores it.

Template errors show the line number in the snippet file, front matter included, with the lines around it.

## Snippet Locations

By default, scaffold looks for snippets in a `snippets` directory next to your scaffoldrc file, for example `~/.config/scaffold/snippets`. Use the [`snippets`](/configuration/scaffold-rc#snippets) key in your scaffoldrc to point at one or more other locations, such as a directory in your dotfiles repo.

```yaml
snippets:
  - ~/dotfiles/snippets
  - ~/work/templates/**/*.snippet.md
```

Each entry is a directory, a single file, or a glob pattern. Directories are searched recursively, symlinks are followed, and hidden files and directories are skipped. Globs support `**` to match any number of directories. A location that can't be read is skipped with a warning.

You can add locations for a single run with the `--snippet-dir` flag or the `SCAFFOLD_SNIPPET_DIR` environment variable. These are searched after the scaffoldrc locations.

A snippet's name is its path relative to the location it was found in, such as `readme.md` or `go/main.go`. You can leave out the file extension when you refer to it. If two locations have a snippet with the same name, the one listed first wins.

## Usage

```bash
# pick a snippet from a searchable list
scaffold snip

# render a snippet by name
scaffold snip readme

# pass variables, using the same syntax as scaffold new
scaffold snip readme name=scaffold description="A scaffolding tool"

# list the available snippets
scaffold snip --list
```

If the name doesn't match a snippet, scaffold suggests the closest match.

## Output

The rendered snippet is written to the current directory with the snippet's file name. `scaffold snip go/main` writes `./main.go`.

Scaffold never replaces an existing file unless you pass `--overwrite`. If the snippet file is executable, the output file is executable too.

Use `--output` (`-o`) to pick another destination:

```bash
# write to a specific file
scaffold snip -o docs/README.md readme

# write into an existing directory, using the snippet's file name
scaffold snip -o docs readme

# write to stdout
scaffold snip -o - gitignore > .gitignore
```

With `-o -`, prompts and messages are written to stderr so that only the rendered snippet goes to stdout.

## Non-Interactive Use

`--no-prompt` skips the questions and messages. Variables come from your scaffoldrc `defaults`, then the `--preset` you pass, then the CLI variables, with later sources taking precedence.

```bash
scaffold snip --no-prompt --preset default readme name=scaffold
```
