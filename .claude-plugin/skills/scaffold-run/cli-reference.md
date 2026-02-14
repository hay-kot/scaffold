# Scaffold CLI Reference

## Commands

### `scaffold new [scaffold] [variables...]`

Generate a project or files from a scaffold template.

| Flag           | Type   | Default | Env Var              | Description                                    |
| -------------- | ------ | ------- | -------------------- | ---------------------------------------------- |
| `--no-prompt`  | bool   | `false` | —                    | Disable interactive mode                       |
| `--preset`     | string | —       | —                    | Preset name for variable values                |
| `--snapshot`   | string | —       | —                    | Path or `stdout` for AST output                |
| `--overwrite`  | bool   | `false` | `SCAFFOLD_OVERWRITE` | Overwrite existing files                       |
| `--force`      | bool   | `true`  | `SCAFFOLD_FORCE`     | Allow dirty git working tree                   |
| `--output-dir` | string | `.`     | `SCAFFOLD_OUT`       | Output directory (`:memory:` for in-memory FS) |
| `--dry-run`    | bool   | `false` | —                    | Validate and show files as JSON                |

### `scaffold list` (alias: `ls`)

List available scaffolds (aliases, local, system).

| Flag     | Type   | Default | Description                    |
| -------- | ------ | ------- | ------------------------------ |
| `--cwd`  | string | `.`     | Working directory to list from |
| `--json` | bool   | `false` | Output JSON                    |

### `scaffold inspect [scaffold]`

Output scaffold metadata as JSON (questions, presets, computed, features, messages).

No subcommand-specific flags. Takes one positional argument.

### `scaffold lint [scaffold.yaml]`

Validate a scaffold.yaml file. Checks:

- Question variable names (alphanumeric + underscore only)
- Prompt types (must be a recognized type)
- Computed variable names
- Skip glob patterns
- Rewrite source paths exist
- Injection modes (`before` or `after`)
- Delimiter globs and values

### `scaffold init`

Create a `.scaffold/` directory with example template scaffold.

| Flag        | Type | Default | Description                            |
| ----------- | ---- | ------- | -------------------------------------- |
| `--stealth` | bool | `false` | Add `.scaffold` to `.git/info/exclude` |

### `scaffold update`

Pull latest versions of all cached system scaffolds.

### `scaffold schema`

Output JSON schema to stdout.

| Flag     | Type   | Default    | Description                             |
| -------- | ------ | ---------- | --------------------------------------- |
| `--type` | string | `scaffold` | Schema type: `scaffold` or `scaffoldrc` |

## Global Flags

| Flag             | Type     | Default           | Env Var(s)                                          | Description                              |
| ---------------- | -------- | ----------------- | --------------------------------------------------- | ---------------------------------------- |
| `--scaffoldrc`   | string   | auto-detected     | `SCAFFOLDRC`                                        | Path to scaffoldrc file                  |
| `--scaffold-dir` | []string | `["./.scaffold"]` | `SCAFFOLD_DIR`                                      | Template directories                     |
| `--cache`        | string   | auto-detected     | `SCAFFOLD_CACHE`                                    | Cache directory                          |
| `--log-level`    | string   | `warn`            | `SCAFFOLD_LOG_LEVEL`, `SCAFFOLD_SETTINGS_LOG_LEVEL` | Log level                                |
| `--log-file`     | string   | —                 | `SCAFFOLD_SETTINGS_LOG_FILE`                        | Log file (`stdout` for stdout)           |
| `--theme`        | string   | `scaffold`        | `SCAFFOLD_SETTINGS_THEME`, `SCAFFOLD_THEME`         | UI theme                                 |
| `--run-hooks`    | string   | `prompt`          | `SCAFFOLD_SETTINGS_RUN_HOOKS`                       | Hook policy: `never`, `always`, `prompt` |

### Theme options

`scaffold` (default), `charm`, `dracula`, `base16`, `catppuccino`, `tokyo-night`

### Hook policies

| Value    | Aliases       | Behavior                     |
| -------- | ------------- | ---------------------------- |
| `prompt` | `""`          | Ask before running (default) |
| `always` | `yes`, `true` | Run without asking           |
| `never`  | `no`, `false` | Never run                    |

## Environment Variables

| Variable                      | Maps to          | Description              |
| ----------------------------- | ---------------- | ------------------------ |
| `SCAFFOLDRC`                  | `--scaffoldrc`   | ScaffoldRC file path     |
| `SCAFFOLD_DIR`                | `--scaffold-dir` | Template directories     |
| `SCAFFOLD_CACHE`              | `--cache`        | Cache directory          |
| `SCAFFOLD_LOG_LEVEL`          | `--log-level`    | Log level                |
| `SCAFFOLD_SETTINGS_LOG_LEVEL` | `--log-level`    | Log level (alt)          |
| `SCAFFOLD_SETTINGS_LOG_FILE`  | `--log-file`     | Log file path            |
| `SCAFFOLD_SETTINGS_THEME`     | `--theme`        | UI theme                 |
| `SCAFFOLD_THEME`              | `--theme`        | UI theme (alt)           |
| `SCAFFOLD_SETTINGS_RUN_HOOKS` | `--run-hooks`    | Hook execution policy    |
| `SCAFFOLD_OVERWRITE`          | `--overwrite`    | Overwrite existing files |
| `SCAFFOLD_FORCE`              | `--force`        | Allow dirty git tree     |
| `SCAFFOLD_OUT`                | `--output-dir`   | Output directory         |

## Scaffold Reference Resolution

When you pass a scaffold name to `scaffold new`, it resolves in this order:

1. **Alias** — exact match against scaffoldrc `aliases` keys
2. **Short** — colon syntax expanded via scaffoldrc `shorts` (e.g., `gh:org/repo`)
3. **URL** — cloned to cache directory
4. **Absolute path** — used directly
5. **Relative path** (contains `/`) — joined with cwd
6. **Bare name** — searched in all `--scaffold-dir` directories

If resolution fails interactively, a fuzzy "did you mean?" prompt is shown.

## Default File Paths

| Path                                       | Purpose    | Legacy Fallback              |
| ------------------------------------------ | ---------- | ---------------------------- |
| `$XDG_CONFIG_HOME/scaffold/scaffoldrc.yml` | ScaffoldRC | `~/.scaffold/scaffoldrc.yml` |
| `$XDG_DATA_HOME/scaffold/templates`        | Cache      | `~/.scaffold/cache`          |

Legacy paths are used if they exist on disk. Use `scaffold dev migrate` to move to XDG locations.
