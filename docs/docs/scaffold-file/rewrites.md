---
title: Rewrites
tags:
  - Scaffold File Reference
---

Rewrites working with the "bootstrap" scaffolds to perform a path rewrite to another directory. The following example defines a rewrite that will render the `templates/defaults.yaml` file to the `roles/{{ .ProjectKebab }}/defaults/main.yaml` path.

```yaml
rewrite:
  - from: templates/defaults.yaml
    to: roles/{{ .ProjectKebab }}/defaults/main.yaml
```

- `from` - The path to the template file
- `to` - a template path to the destination file
- These files _are_ rendered with the template engine

This feature is not available for project scaffolds.