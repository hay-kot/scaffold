---
description: Lint, dry-run, and snapshot test a scaffold
allowed-tools:
  - Bash
  - Read
  - Edit
  - Glob
argument-hint: "<scaffold-path> [preset]"
---

# Test a scaffold

Validate and test a scaffold template using the scaffold CLI's built-in testing tools.

## Workflow

### 1. Identify the scaffold

If `$ARGUMENTS` provides a path, use it. Otherwise, list available scaffolds and ask:

```bash
scaffold list --json
```

### 2. Lint

Validate the scaffold.yaml configuration:

```bash
scaffold lint <scaffold-path>/scaffold.yaml
```

Report any validation errors. If lint fails, help fix the issues before proceeding.

### 3. Inspect

Check scaffold metadata:

```bash
scaffold inspect <scaffold-path>
```

Verify questions, presets, computed values, and features look correct.

### 4. Dry-run

Run without writing to disk:

```bash
scaffold new --dry-run --no-prompt --preset <preset> <scaffold-path>
```

Use `default` preset unless `$2` specifies otherwise. Report the files that would be created and any errors.

### 5. Snapshot test

Full in-memory render with AST output:

```bash
scaffold --log-level error new --output-dir :memory: --no-prompt --preset <preset> --snapshot stdout <scaffold-path>
```

Review the rendered output for correctness. Check:

- All expected files are present
- Template variables resolved correctly
- Conditional content rendered as expected
- No template errors in output

### 6. Test all presets

If the scaffold has multiple presets, test each one:

```bash
scaffold inspect <scaffold-path>
```

Then run the snapshot test for each preset found.

### 7. Report

Summarize results:

- Lint: pass/fail
- Dry-run: files that would be created
- Snapshot: any issues found in rendered output
- Per-preset results (if multiple presets)
