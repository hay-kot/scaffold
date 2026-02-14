# scaffold-author — Create, edit, and test scaffold templates

User-invocable: true
Auto-triggerable: true

Triggers when: User wants to create a new scaffold template, edit scaffold.yaml, write template files, or test/validate a scaffold.

## Scaffold Types

### Project Scaffold

Creates a new directory named after the project. The user is prompted for a project name.

```
my-scaffold/
├── scaffold.yaml
└── {{ .ProjectKebab }}/          # or any project name pattern
    ├── README.md
    └── src/
        └── main.go
```

Valid root directory names: `{{ .Project }}`, `{{ .ProjectSlug }}`, `{{ .ProjectSnake }}`, `{{ .ProjectKebab }}`, `{{ .ProjectCamel }}`, `{{ .ProjectPascal }}`

### Template Scaffold

Generates files into the current directory. No project name prompt. The `templates/` prefix is stripped from output paths. Use `rewrites` and `inject` to place files precisely.

```
my-scaffold/
├── scaffold.yaml
└── templates/
    ├── component.tsx
    └── component.test.tsx
```

## Quick Start: Project Scaffold

```bash
scaffold init  # creates .scaffold/ with example template scaffold
```

Or create manually:

```
.scaffold/my-project/
├── scaffold.yaml
└── {{ .ProjectKebab }}/
    └── main.go
```

**scaffold.yaml:**

```yaml
questions:
  - name: description
    prompt:
      message: "Project description"
    validate:
      required: true

  - name: license
    prompt:
      message: "License"
      default: "MIT"
      options: ["MIT", "Apache-2.0", "GPL-3.0"]

presets:
  default:
    description: "My project"
    license: "MIT"
```

**`{{ .ProjectKebab }}/main.go`:**

```
// {{ .Scaffold.description }}
package main

func main() {}
```

## Quick Start: Template Scaffold

```
.scaffold/add-component/
├── scaffold.yaml
└── templates/
    └── component.tsx
```

**scaffold.yaml:**

```yaml
questions:
  - name: component_name
    prompt:
      message: "Component name"
    validate:
      required: true

rewrites:
  - from: templates/component.tsx
    to: "src/components/{{ .Scaffold.component_name }}.tsx"

presets:
  default:
    component_name: "Button"
```

## Bootstrapping

### `scaffold init`

Creates `.scaffold/` with a hello-world template scaffold example.

### `scaffold init --stealth`

Same as above, but adds `.scaffold` to `.git/info/exclude` so it's hidden from git without touching `.gitignore`. Useful for personal scaffolds in shared repos.

## scaffold.yaml Reference

See `scaffold-yaml-schema.md` for the complete field reference. Key sections:

- **questions** — Interactive prompts that collect user input
- **computed** — Derived variables from template expressions
- **features** — Conditional file inclusion/exclusion
- **rewrites** — Remap output file paths (template scaffolds)
- **inject** — Insert content into existing files (template scaffolds)
- **presets** — Named variable sets for non-interactive use
- **each** — Multi-file expansion from list variables
- **skip** — Glob patterns for files copied without template rendering
- **delimiters** — Custom template delimiters per file glob
- **messages** — Pre/post scaffold markdown messages
- **metadata** — Version constraints

## Template Engine

See `template-engine.md` for the complete template reference. Key points:

- Variables: `.Project`, `.ProjectKebab`, `.Scaffold.*`, `.Computed.*`, `.Each.*`
- Functions: All sprout functions + `wraptmpl`, `toPlural`, `toSingular`, `isPlural`, `isSingular`, `partial`
- File/directory names are also templates
- Partials go in a `partials/` directory at the scaffold root

## Testing Workflow

### 1. Lint the config

```bash
scaffold lint .scaffold/my-scaffold/scaffold.yaml
```

Validates question names, prompt types, computed names, skip globs, rewrite paths, injection modes, and delimiter configs.

### 2. Inspect the metadata

```bash
scaffold inspect .scaffold/my-scaffold
```

Outputs JSON showing questions (with resolved types), presets, computed values, features, and messages. Use this to verify the scaffold looks correct.

### 3. Dry-run

```bash
scaffold new --dry-run --no-prompt --preset default .scaffold/my-scaffold
```

Renders fully but writes nothing. Outputs JSON listing all files that would be created. Catches template errors without side effects.

### 4. Snapshot test

```bash
scaffold new \
  --output-dir :memory: \
  --no-prompt \
  --preset default \
  --snapshot stdout \
  .scaffold/my-scaffold
```

Renders in-memory and outputs a full AST with file contents. Ideal for:

- Verifying exact output content
- CI/CD diffing against expected output
- Catching regressions

### 5. Full test command (recommended)

```bash
scaffold \
  --log-level error \
  new \
  --output-dir :memory: \
  --no-prompt \
  --preset default \
  --snapshot stdout \
  .scaffold/my-scaffold
```

## Reading Template Errors

Template errors show:

- **File path** and **line:column** where the error occurred
- **Context lines** (the error line plus surrounding lines)
- **Delimiter info** if custom delimiters are in use

Common causes:

- Referencing undefined variables (e.g., `.Scaffold.typo`)
- Missing closing delimiters
- Calling undefined functions
- Type mismatches in template expressions

## Preset-Driven Testing

Always define at least one preset (typically `default`) so scaffolds can be tested non-interactively. The preset should provide valid values for all required questions.

```yaml
presets:
  default:
    description: "Test project"
    license: "MIT"
    use_database: true
    db_type: "postgres"
  minimal:
    description: "Minimal"
    license: "MIT"
    use_database: false
```

Test each preset:

```bash
for preset in default minimal; do
  scaffold new --output-dir :memory: --no-prompt --preset "$preset" --snapshot stdout .scaffold/my-scaffold
done
```
