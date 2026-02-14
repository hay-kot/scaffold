# Template Engine Reference

The scaffold template engine is built on Go's `text/template` with sprout functions and custom additions.

## Template Variables

All variables are available in file contents and file/directory name templates.

### Top-Level Variables

| Variable         | Type   | Description             |
| ---------------- | ------ | ----------------------- |
| `.Project`       | string | Project name as entered |
| `.ProjectSnake`  | string | `snake_case`            |
| `.ProjectKebab`  | string | `kebab-case`            |
| `.ProjectCamel`  | string | `camelCase`             |
| `.ProjectPascal` | string | `PascalCase`            |
| `.ProjectSlug`   | string | URL-safe slug           |

### `.Scaffold.*`

User-provided question answers, keyed by question `name`.

```
{{ .Scaffold.description }}     → string
{{ .Scaffold.use_docker }}      → bool
{{ .Scaffold.features }}        → []string
```

The type matches the question prompt type: string (text/select), bool (confirm), []string (multi-select/loop).

### `.Computed.*`

Derived variables defined in scaffold.yaml `computed` section.

```
{{ .Computed.has_auth }}
{{ .Computed.project_upper }}
```

Computed values undergo type coercion: rendered string is parsed as int, then bool, then kept as string.

### `.Each.*`

Available only inside `each`-expanded templates.

| Variable      | Type   | Description        |
| ------------- | ------ | ------------------ |
| `.Each.Item`  | string | Current item value |
| `.Each.Index` | int    | Zero-based index   |

## Template Functions

### Sprout Functions

All [sprout](https://docs.atom.codes/sprout) functions are available. Common ones:

**String:**
`upper`, `lower`, `trim`, `trimSuffix`, `trimPrefix`, `replace`, `contains`, `hasPrefix`, `hasSuffix`, `repeat`, `nospace`, `indent`, `nindent`, `quote`, `squote`, `title`, `untitle`, `camelcase`, `snakecase`, `kebabcase`, `shuffle`, `substr`, `trunc`

**Collections:**
`list`, `append`, `prepend`, `first`, `last`, `has`, `without`, `uniq`, `join`, `sortAlpha`, `reverse`, `concat`, `chunk`

**Math:**
`add`, `sub`, `mul`, `div`, `mod`, `max`, `min`, `ceil`, `floor`, `round`

**Type conversion:**
`toString`, `toInt`, `toFloat64`, `toBool`

**Conditionals:**
`ternary`, `default`, `empty`, `coalesce`

**Flow:**
`fail`

### Custom Functions

| Function     | Signature                               | Description                                                                |
| ------------ | --------------------------------------- | -------------------------------------------------------------------------- |
| `wraptmpl`   | `wraptmpl(s string) string`             | Wraps string in `{{ }}`. Use when output should contain Go template syntax |
| `toPlural`   | `toPlural(s string) string`             | Convert word to plural form                                                |
| `toSingular` | `toSingular(s string) string`           | Convert word to singular form                                              |
| `isPlural`   | `isPlural(s string) bool`               | Check if word is plural                                                    |
| `isSingular` | `isSingular(s string) bool`             | Check if word is singular                                                  |
| `partial`    | `partial(name string, data any) string` | Render a named partial with given data context                             |

### `wraptmpl` example

When your scaffold output should contain literal Go template syntax:

```
{{ wraptmpl ".Values.image" }}
```

Produces: `{{ .Values.image }}`

Alternative approach using raw strings:

```
{{ "{{ .Values.image }}" }}
```

## File and Directory Name Templating

File and directory names are templates too:

```
{{ .ProjectKebab }}/
├── {{ .Scaffold.module_name }}.go
├── {{ .Scaffold.module_name }}_test.go
└── cmd/
    └── {{ .ProjectSnake }}/
        └── main.go
```

Names without `{{ }}` syntax are passed through unchanged (no parsing overhead).

## Partials

Reusable template fragments stored in a `partials/` directory at the scaffold root.

### Directory structure

```
my-scaffold/
├── scaffold.yaml
├── partials/
│   ├── header.txt
│   └── license.txt
└── {{ .ProjectKebab }}/
    └── README.md
```

### Registration rules

- File extension is stripped from the partial name
- Names must contain only letters, digits, and underscores
- `partials/header.txt` → partial name `header`
- `partials/license.txt` → partial name `license`

### Usage in templates

```
{{ partial "header" . }}

# {{ .Project }}

{{ partial "license" . }}
```

The second argument (`.`) passes the full variable context to the partial. Partials have access to all the same variables as regular templates.

## Engine Rules

1. **Empty source files are skipped** — files with zero bytes are not processed
2. **Templates rendering to empty output are excluded** — if the rendered output is empty or whitespace-only, the file is not written
3. **Feature flags filter files** — files matching a feature's globs are excluded when the feature's value evaluates to `false`
4. **Skip patterns bypass rendering** — files matching `skip` globs are copied as-is without template processing
5. **Custom delimiters apply per-file** — the first matching delimiter glob wins
6. **Guard chain order**: rewrite → render path → no-clobber check → directory handling → feature flag check

## Common Patterns

### Conditional content blocks

```
{{ if .Scaffold.use_docker }}
COPY . /app
RUN go build -o /app/main .
{{ end }}
```

### Iterating over lists

```
{{ range .Scaffold.features }}
- {{ . }}
{{ end }}
```

### Conditional file inclusion

Use `features` in scaffold.yaml instead of wrapping entire file contents in `{{ if }}`. Feature-excluded files produce no output file at all.

### Preserving Go template syntax in output

Three approaches:

1. `{{ wraptmpl ".Values.name" }}` — custom function
2. `{{ "{{ .Values.name }}" }}` — raw string
3. Add to `skip` list — bypasses template engine entirely
4. Use custom `delimiters` — use different syntax for scaffold vs output templates
