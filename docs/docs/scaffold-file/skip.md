---
title: Skip
tags:
  - Scaffold File Reference
---

Skip is a list of glob patterns that will be used to skip the template **rendering** process. This is useful is your file is itself a go template, or contains similar syntax that will cause the template engine to fail. The following example will skip the `templates/defaults.yaml` file from being rendered.

```yaml
skip:
  - "*.goreleaser.yaml"
  - "**/*.gotmpl"
```