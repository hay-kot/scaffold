---
description: Create a new scaffold template in the project
allowed-tools:
  - Bash
  - Read
  - Write
  - Edit
  - Glob
  - Grep
argument-hint: "[scaffold-name]"
---

# Create a scaffold template

Help the user create a new scaffold template in the project's `.scaffold/` directory.

## Workflow

### 1. Determine scaffold type

Ask the user:

- **Project scaffold**: Creates a new directory named after the project. Root directory uses a project name variable like `{{ .ProjectKebab }}`.
- **Template scaffold**: Generates files into the current directory using a `templates/` prefix. Uses `rewrites` to place files.

### 2. Initialize if needed

If `.scaffold/` doesn't exist:

```bash
scaffold init
```

### 3. Create scaffold structure

If `$ARGUMENTS` provides a name, use it. Otherwise ask the user for a scaffold name.

**Project scaffold:**

```
.scaffold/<name>/
├── scaffold.yaml
└── {{ .ProjectKebab }}/
    └── (template files)
```

**Template scaffold:**

```
.scaffold/<name>/
├── scaffold.yaml
└── templates/
    └── (template files)
```

### 4. Define scaffold.yaml

Work with the user to define:

- **questions** — Interactive prompts for user input
- **presets** — At least one `default` preset for testing
- **features** — Conditional file inclusion (if needed)
- **rewrites** — Output path remapping (template scaffolds)
- **inject** — Content injection into existing files (if needed)
- **computed** — Derived variables (if needed)

Refer to the scaffold-author skill for the complete scaffold.yaml schema and template engine reference.

### 5. Create template files

Write the template files using Go template syntax. Available variables:

- `.Project`, `.ProjectKebab`, `.ProjectSnake`, `.ProjectCamel`, `.ProjectPascal`
- `.Scaffold.<question_name>` for user answers
- `.Computed.<name>` for computed values

### 6. Validate

After creating the scaffold, run validation:

```bash
scaffold lint .scaffold/<name>/scaffold.yaml
scaffold new --output-dir :memory: --no-prompt --preset default --snapshot stdout .scaffold/<name>
```
