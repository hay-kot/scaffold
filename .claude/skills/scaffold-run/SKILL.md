# scaffold-run — Use scaffold CLI to generate projects/files

User-invocable: true
Auto-triggerable: true

Triggers when: User wants to run a scaffold, generate a project, create files from a template, or use the scaffold CLI.

## Workflow

### 1. Discover available scaffolds

```bash
scaffold list --json
```

Lists all scaffolds: aliases (from scaffoldrc), local (in `.scaffold/`), and system (cached remote repos). Use `--json` for programmatic output.

### 2. Inspect a scaffold before running

```bash
scaffold inspect <name-or-path>
```

Outputs JSON with questions (name, type, default, options), presets, computed values, features, and messages. Use this to understand what variables a scaffold expects.

### 3. Run a scaffold interactively

```bash
scaffold new <scaffold>
```

Launches an interactive prompt for project name and all questions. The scaffold argument can be:
- An alias name (defined in scaffoldrc)
- A short reference (`gh:org/repo`)
- A URL (`https://github.com/org/repo`)
- An absolute or relative filesystem path
- A bare name (searched in `--scaffold-dir` directories, default `.scaffold/`)

### 4. Run non-interactively

```bash
scaffold new --no-prompt --preset <preset-name> <scaffold> [key=value ...]
```

Use `--no-prompt` to skip all interactive prompts. Provide variables via:
- `--preset <name>` — loads a named preset from the scaffold's config
- Positional `key[:type]=value` arguments — override or supplement preset values

CLI arguments take precedence over preset values.

If no `Project` variable is set in non-interactive mode, a random name `scaffold-test-NNNN` is auto-generated.

### 5. Validate with dry-run

```bash
scaffold new --dry-run --no-prompt --preset default <scaffold>
```

Renders the scaffold fully but writes nothing to disk. Outputs JSON:

```json
{
  "files": [{"path": "path/to/file", "action": "create"}],
  "errors": [],
  "warnings": []
}
```

### 6. In-memory testing with snapshot

```bash
scaffold new --output-dir :memory: --no-prompt --preset default --snapshot stdout <scaffold>
```

Renders entirely in memory and outputs a full AST with file contents to stdout. Ideal for CI/CD validation and diffing.

### 7. Passing typed variables

Variables use `key[:type]=value` syntax. See `variable-syntax.md` for the full type reference.

Common examples:
```bash
scaffold new --no-prompt my-scaffold \
  Project=MyApp \
  description="A web service" \
  port:int=8080 \
  debug:bool=true \
  features:[]string=auth,api
```

## Key flags for `scaffold new`

| Flag | Default | Description |
|------|---------|-------------|
| `--no-prompt` | `false` | Disable interactive prompts |
| `--preset` | — | Preset name for variable values |
| `--output-dir` | `.` | Output directory (`:memory:` for in-memory) |
| `--dry-run` | `false` | Show what would be created (JSON) |
| `--snapshot` | — | Output AST to path or `stdout` |
| `--overwrite` | `false` | Overwrite existing files |
| `--force` | `true` | Allow dirty git working tree |

## Tips

- Always `scaffold inspect` first to understand a scaffold's variables and presets
- Use `--dry-run` before real runs to verify output paths
- Combine `--output-dir :memory:` with `--snapshot stdout` for zero-disk testing
- When `--no-prompt` is active, hooks with `run_hooks=prompt` are skipped
- Pre/post messages are suppressed in `--no-prompt` mode
