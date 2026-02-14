# scaffold.yaml Schema Reference

## Top-Level Structure

```yaml
metadata:
  minimum_version: "1.5.0"

messages:
  pre: "..."
  post: "..."

questions: [...]
computed: { ... }
features: [...]
presets: { ... }
each: [...]
skip: [...]
delimiters: [...]
rewrites: [...] # template scaffolds only
inject: [...] # template scaffolds only
```

---

## `metadata`

```yaml
metadata:
  minimum_version: "1.5.0"
```

| Field             | Type   | Description                                                                     |
| ----------------- | ------ | ------------------------------------------------------------------------------- |
| `minimum_version` | string | Minimum scaffold CLI version required (semver). Empty or `"*"` skips the check. |

---

## `questions`

Interactive prompts that collect user input. Each question's answer is available as `{{ .Scaffold.<name> }}` in templates.

```yaml
questions:
  - name: "variable_name"
    group: "group_name"
    when: "{{ .previous_var }}"
    required: true
    validate:
      required: true
      min: 3
      max: 50
      match:
        regex: "^[a-z][a-z0-9-]*$"
        message: "Lowercase letters, numbers, hyphens only"
    prompt:
      message: "Display label"
      description: "Help text"
      default: "default value"
      options: ["a", "b", "c"]
      multi: false
      loop: false
      confirm: "Yes or no?"
```

### Question fields

| Field      | Type   | Description                                                                                                                 |
| ---------- | ------ | --------------------------------------------------------------------------------------------------------------------------- |
| `name`     | string | Variable name (alphanumeric + underscore). Accessed as `.Scaffold.<name>`                                                   |
| `group`    | string | Groups questions into a single UI form. When used with `when`, only the first question's `when` applies to the entire group |
| `when`     | string | Go template evaluated with previous answers at root level. If `false`, question is skipped                                  |
| `required` | bool   | Whether input is required                                                                                                   |
| `validate` | object | Validation config (see below)                                                                                               |
| `prompt`   | object | Prompt display config (see below)                                                                                           |

### Prompt types

The prompt type is determined by which fields are set:

| Fields Set                             | Type              | Result Type |
| -------------------------------------- | ----------------- | ----------- |
| `message` + `options` + `multi: true`  | Multi Select      | `[]string`  |
| `message` + `options`                  | Select            | `string`    |
| `confirm`                              | Confirm           | `bool`      |
| `message` + `loop: true`               | Looped Text Input | `[]string`  |
| `message` + `multi: true` (no options) | Multiline Text    | `string`    |
| `message` only                         | Text Input        | `string`    |

### Prompt fields

| Field         | Type     | Description                                                     |
| ------------- | -------- | --------------------------------------------------------------- |
| `message`     | string   | Input label                                                     |
| `description` | string   | Help text below input                                           |
| `default`     | any      | Default value (string, bool, or []string depending on type)     |
| `confirm`     | string   | Creates a yes/no prompt with this as the title                  |
| `options`     | []string | Creates a select input (or multi-select with `multi: true`)     |
| `multi`       | bool     | Multi-select (with options) or multiline text (without options) |
| `loop`        | bool     | Looped text input producing `[]string`                          |

### Validate fields

| Field           | Type   | Description                                                                        |
| --------------- | ------ | ---------------------------------------------------------------------------------- |
| `required`      | bool   | Input must not be empty                                                            |
| `min`           | int    | Minimum length (chars for strings, items for slices). If > 0, overrides `required` |
| `max`           | int    | Maximum length                                                                     |
| `match.regex`   | string | Regex the input must match. For slices, each element is validated                  |
| `match.message` | string | Error message when regex doesn't match                                             |

### Prompt examples

**Text input:**

```yaml
- name: description
  prompt:
    message: "Project description"
```

**Select:**

```yaml
- name: license
  prompt:
    message: "License"
    default: "MIT"
    options: ["MIT", "Apache-2.0", "GPL-3.0"]
```

**Multi-select:**

```yaml
- name: features
  prompt:
    message: "Features to enable"
    multi: true
    default: ["auth"]
    options: ["auth", "api", "worker", "cron"]
```

**Confirm:**

```yaml
- name: use_docker
  prompt:
    confirm: "Include Docker support?"
```

**Looped input (produces []string):**

```yaml
- name: services
  prompt:
    message: "Service name"
    description: "Enter service names one at a time"
    loop: true
  validate:
    min: 1
```

**Multiline text:**

```yaml
- name: readme_content
  prompt:
    message: "README content"
    multi: true
```

**Conditional question:**

```yaml
- name: db_type
  when: "{{ .use_docker }}"
  prompt:
    message: "Database type"
    options: ["postgres", "mysql", "sqlite"]
```

---

## `computed`

Derived variables from Go template expressions. Available as `{{ .Computed.<name> }}`.

```yaml
computed:
  project_upper: "{{ .Project | upper }}"
  has_auth: '{{ has "auth" .Scaffold.features }}'
```

| Key           | Value                                           |
| ------------- | ----------------------------------------------- |
| variable name | Go template string evaluated with all variables |

Computed values undergo type coercion: the rendered string is parsed as int first, then bool, then kept as string. This means `"0"` becomes `0` (int), `"true"` becomes `true` (bool).

---

## `features`

Conditionally include or exclude files based on template expressions.

```yaml
features:
  - value: "{{ .Scaffold.use_docker }}"
    globs:
      - "docker/**/*"
      - "Dockerfile"
  - value: "{{ .Scaffold.use_ci }}"
    globs:
      - ".github/**/*"
```

| Field   | Type     | Description                                                                      |
| ------- | -------- | -------------------------------------------------------------------------------- |
| `value` | string   | Go template. If result is `true`, matched files are included; otherwise excluded |
| `globs` | []string | Glob patterns to match against template file paths                               |

---

## `presets`

Named sets of variable values for non-interactive execution.

```yaml
presets:
  default:
    description: "My project"
    license: "MIT"
    use_docker: true
    features: ["auth", "api"]
  minimal:
    description: "Minimal"
    license: "MIT"
    use_docker: false
    features: []
```

Each preset maps question names to values. Used with `scaffold new --preset <name>`.

---

## `each`

Multi-file expansion from list variables. Files or directories with `[varname]` in their path are expanded once per item in the list.

### String shorthand

```yaml
each:
  - services
```

### Object form with path transformation

```yaml
each:
  - var: models
    as: "{{ .Each.Item | toPascalCase }}"
```

| Field | Type   | Description                                                                              |
| ----- | ------ | ---------------------------------------------------------------------------------------- |
| `var` | string | Name of the list variable (must produce `[]string`)                                      |
| `as`  | string | Optional Go template for path segment transformation. Has `.Each.Item` and `.Each.Index` |

Inside expanded templates, use:

- `{{ .Each.Item }}` — current string item
- `{{ .Each.Index }}` — zero-based index

**Directory expansion:** `[services]/` — every file inside is rendered per item.
**File expansion:** `[services].go` — that file is rendered per item.

Undeclared `[varname]` segments are treated as literal path characters.

---

## `skip`

Glob patterns for files that should be copied without template rendering. The file contents are preserved as-is; path guards (rewrites, feature flags) still apply.

```yaml
skip:
  - "*.goreleaser.yaml"
  - "**/*.gotmpl"
  - ".github/workflows/*.yml"
```

Use this for files that contain Go template syntax (`{{ }}`) that should be preserved literally.

---

## `delimiters`

Custom template delimiters for files that conflict with the default `{{ }}` syntax.

```yaml
delimiters:
  - glob: "**.goreleaser"
    left: "[["
    right: "]]"
  - glob: "**/*.helmfile"
    left: "<%"
    right: "%>"
```

| Field   | Type   | Description                       |
| ------- | ------ | --------------------------------- |
| `glob`  | string | File pattern to match             |
| `left`  | string | Opening delimiter (replaces `{{`) |
| `right` | string | Closing delimiter (replaces `}}`) |

The first matching glob wins for each file.

---

## `rewrites`

Remap output file paths. **Template scaffolds only.**

```yaml
rewrites:
  - from: templates/config.yaml
    to: "config/{{ .Scaffold.app_name }}.yaml"
  - from: templates/defaults.yaml
    to: "roles/{{ .ProjectKebab }}/defaults/main.yaml"
```

| Field  | Type   | Description                             |
| ------ | ------ | --------------------------------------- |
| `from` | string | Source path (relative to scaffold root) |
| `to`   | string | Go template for output path             |

---

## `inject`

Insert content into existing files. **Template scaffolds only.**

```yaml
inject:
  - name: "register route"
    path: src/routes.ts
    at: "// ROUTES"
    mode: after
    template: |
      app.use("/{{ .Scaffold.route_name }}", {{ .Scaffold.route_name }}Router);
```

| Field      | Type   | Description                                                                      |
| ---------- | ------ | -------------------------------------------------------------------------------- |
| `name`     | string | Descriptive name                                                                 |
| `path`     | string | Target file path (may be a template)                                             |
| `at`       | string | Substring to locate injection point. ALL matches are replaced                    |
| `mode`     | string | `before` (default) or `after` the matched line                                   |
| `template` | string | Go template to inject. If empty/whitespace after rendering, injection is skipped |

The injected content inherits the indentation of the matched line.

---

## `messages`

Markdown messages shown before and after scaffold generation.

```yaml
messages:
  pre: |
    # Welcome
    This scaffold creates a new Go service.
    Template variables are **not** available here.

  post: |
    # Done
    Your project **{{ .ProjectKebab }}** is ready.
    Run `cd {{ .ProjectKebab }} && go run .` to start.
```

| Field  | Type   | Description                                               |
| ------ | ------ | --------------------------------------------------------- |
| `pre`  | string | Shown before generation. Template variables NOT available |
| `post` | string | Shown after generation. Template variables ARE available  |

Both are rendered as markdown in the terminal. Suppressed in `--no-prompt` mode.
