---
description: Discover and run a scaffold to generate a project or files
allowed-tools:
  - Bash
  - Read
  - Glob
argument-hint: "[scaffold-name] [key=value ...]"
---

# Run a scaffold

Help the user discover and run a scaffold using the scaffold CLI.

## Workflow

### 1. Identify the scaffold

If `$ARGUMENTS` specifies a scaffold name, use it directly. Otherwise, discover available scaffolds:

```bash
scaffold list --json
```

Present the available scaffolds to the user and ask which one to run.

### 2. Inspect the scaffold

Before running, inspect the scaffold to understand its questions and presets:

```bash
scaffold inspect <scaffold>
```

Show the user what variables the scaffold expects.

### 3. Run the scaffold

Run the scaffold interactively:

```bash
scaffold new <scaffold>
```

If the user provides key=value pairs in `$ARGUMENTS`, pass them through:

```bash
scaffold new <scaffold> $ARGUMENTS
```

For non-interactive execution with a preset:

```bash
scaffold new --no-prompt --preset <preset> <scaffold> [key=value ...]
```

### 4. Verify output

After generation, check the output and summarize what was created. If any errors occurred, help the user debug using the error output.
