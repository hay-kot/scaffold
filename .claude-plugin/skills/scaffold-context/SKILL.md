# scaffold-context — Runtime configuration background knowledge

User-invocable: false
Auto-triggerable: true

Triggers when: User asks about scaffoldrc setup, scaffold resolution, environment variables, authentication, or how scaffold finds templates.

## ScaffoldRC Location

The scaffoldrc file is resolved in this order:

1. `--scaffoldrc` CLI flag
2. `SCAFFOLDRC` environment variable
3. `~/.scaffold/scaffoldrc.yml` (legacy, only if it exists on disk)
4. `$XDG_CONFIG_HOME/scaffold/scaffoldrc.yml` (default: `~/.config/scaffold/scaffoldrc.yml`)

If the file doesn't exist, it's created automatically as an empty file. An empty file is valid and uses all defaults.

Use `scaffold dev migrate` to move legacy `~/.scaffold/` paths to XDG locations.

## Scaffold Reference Resolution

When `scaffold new <ref>` is called, the reference resolves through:

1. **Aliases** — exact match against scaffoldrc `aliases` keys → replaced with alias value
2. **Shorts** — colon syntax (e.g., `gh:org/repo`) expanded via scaffoldrc `shorts` → URL
3. **Remote URL** — cloned to cache directory, specific version/branch/tag supported
4. **Absolute path** — used directly
5. **Relative path** (contains `/`) — joined with current working directory
6. **Bare name** — searched in all `--scaffold-dir` directories (default: `.scaffold/`)

If resolution fails interactively, fuzzy matching runs against all known scaffolds with a "did you mean?" prompt. In `--no-prompt` mode, the error is returned directly.

## Authentication

For private repositories, configure auth in scaffoldrc. If no auth is configured and a remote requires it, the CLI prompts interactively (unless `--no-prompt`).

See `scaffoldrc-reference.md` for the `auth` section schema.

## Environment Variables

| Variable                      | Maps to          | Description                          |
| ----------------------------- | ---------------- | ------------------------------------ |
| `SCAFFOLDRC`                  | `--scaffoldrc`   | ScaffoldRC file path                 |
| `SCAFFOLD_DIR`                | `--scaffold-dir` | Template search directories          |
| `SCAFFOLD_CACHE`              | `--cache`        | Cache directory for remote scaffolds |
| `SCAFFOLD_LOG_LEVEL`          | `--log-level`    | Log level (debug, info, warn, error) |
| `SCAFFOLD_SETTINGS_LOG_LEVEL` | `--log-level`    | Log level (alt)                      |
| `SCAFFOLD_SETTINGS_LOG_FILE`  | `--log-file`     | Log file path                        |
| `SCAFFOLD_SETTINGS_THEME`     | `--theme`        | UI theme                             |
| `SCAFFOLD_THEME`              | `--theme`        | UI theme (alt)                       |
| `SCAFFOLD_SETTINGS_RUN_HOOKS` | `--run-hooks`    | Hook execution policy                |
| `SCAFFOLD_OVERWRITE`          | `--overwrite`    | Overwrite existing files             |
| `SCAFFOLD_FORCE`              | `--force`        | Allow dirty git tree                 |
| `SCAFFOLD_OUT`                | `--output-dir`   | Output directory                     |

CLI flags take precedence over environment variables, which take precedence over scaffoldrc settings.

## Default Paths

| Purpose    | XDG Path                                   | Legacy Path                  |
| ---------- | ------------------------------------------ | ---------------------------- |
| ScaffoldRC | `$XDG_CONFIG_HOME/scaffold/scaffoldrc.yml` | `~/.scaffold/scaffoldrc.yml` |
| Cache      | `$XDG_DATA_HOME/scaffold/templates`        | `~/.scaffold/cache`          |

Legacy paths are used if they exist. XDG defaults: `~/.config/` for config, `~/.local/share/` for data.

## Cache Behavior

Remote scaffolds (URLs, shorts) are cloned to the cache directory. Use `scaffold update` to pull latest for all cached scaffolds. The cache directory can be overridden with `--cache` or `SCAFFOLD_CACHE`.
