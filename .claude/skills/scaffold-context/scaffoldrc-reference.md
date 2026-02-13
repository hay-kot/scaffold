# ScaffoldRC Reference

The scaffoldrc file (`scaffoldrc.yml`) configures global settings, defaults, aliases, shorts, and authentication.

## Full Example

```yaml
settings:
  theme: scaffold
  run_hooks: prompt
  log_level: warn
  log_file: stdout

defaults:
  author: "Jane Doe"
  github_username: "janedoe"
  license: "MIT"

aliases:
  api: ~/scaffolds/go-api
  component: https://github.com/org/scaffold-component

shorts:
  gh: https://github.com/
  gl: https://gitlab.com/

auth:
  - match: "^https://github.com"
    token: ${GITHUB_TOKEN}
  - match: "^https://gitea.example.com"
    basic:
      username: deploy
      password: ${GITEA_PASSWORD}
```

---

## `settings`

| Setting | Type | Default | CLI Flag | Description |
|---------|------|---------|----------|-------------|
| `theme` | string | `scaffold` | `--theme` | UI theme for interactive prompts |
| `run_hooks` | string | `prompt` | `--run-hooks` | Hook execution policy |
| `log_level` | string | `warn` | `--log-level` | Log level |
| `log_file` | string | `stdout` | `--log-file` | Log output destination |

**Theme options:** `scaffold`, `charm`, `dracula`, `base16`, `catppuccino`, `tokyo-night`

**Run hooks options:**

| Value | Aliases | Behavior |
|-------|---------|----------|
| `prompt` | `""` | Ask before running (default) |
| `always` | `yes`, `true` | Run without asking |
| `never` | `no`, `false` | Never run |

CLI flags override scaffoldrc settings when explicitly set.

---

## `defaults`

```yaml
defaults:
  author: "Jane Doe"
  github_username: "janedoe"
```

A map of variable names to default values. Injected into every scaffold run. CLI arguments override defaults (CLI args take precedence during merge).

---

## `aliases`

```yaml
aliases:
  api: ~/scaffolds/go-api
  component: https://github.com/org/scaffold-component
  role: /absolute/path/to/scaffold
```

Maps short names to full scaffold paths. Used as `scaffold new api` instead of the full path.

**Rules:**
- Values must be a valid URL, an absolute path, or start with `~`
- Relative paths without `~` are rejected during validation
- Aliases are shown in the interactive scaffold picker
- Exact match only (no fuzzy matching on alias names)

---

## `shorts`

```yaml
shorts:
  gh: https://github.com/
  gl: https://gitlab.com/
  bb: https://bitbucket.org/
```

URL prefix shortcuts. Used with colon syntax: `gh:org/repo` expands to `https://github.com/org/repo`.

**Rules:**
- Values must be valid URIs
- The input is split on `:` — left side matches the short key, right side is appended via `url.JoinPath`
- Only the first `:` is used for splitting

---

## `auth`

```yaml
auth:
  - match: "^https://github.com"
    token: ${GITHUB_TOKEN}

  - match: "^https://gitea.example.com"
    basic:
      username: deploy
      password: ${GITEA_PASSWORD}
```

Pattern-matched authentication for private git repositories.

| Field | Type | Description |
|-------|------|-------------|
| `match` | regex | Regular expression matched against the scaffold URL |
| `token` | string | Token authentication |
| `basic.username` | string | HTTP basic auth username |
| `basic.password` | string | HTTP basic auth password |

**Rules:**
- Entries are checked in order; first match wins
- Use `token` OR `basic`, not both
- Environment variable expansion: `${VAR_NAME}` is replaced with the env var value
- Only whole-value references are supported (e.g., `${TOKEN}`), not partial interpolation (e.g., `prefix-${TOKEN}` won't work)
- If no auth entry matches and the remote requires authentication, the CLI prompts interactively (unless `--no-prompt`)

---

## JSON Schema

Generate the scaffoldrc JSON schema:

```bash
scaffold schema --type scaffoldrc
```
