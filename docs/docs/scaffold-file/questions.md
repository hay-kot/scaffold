---
title: Interactive Prompts
tags:
  - Scaffold File Reference
---

Questions are used to prompt the user for input when generating a scaffold. There are several types of questions that are available.

Note: You can use the `required` flag to make a question required.

### Text

```yaml
questions:
  - name: "Description"
    prompt:
      message: "Description of the project"
    required: true
```

### Boolean (Yes/No)

```yaml
questions:
  - name: "Use Github Actions"
    prompt:
      confirm: "Use Github Actions for CI/CD?"
```

### Multi Select

```yaml
questions:
  - name: "Colors"
    prompt:
      multi: true
      message: "Colors of the project"
      options:
        - "red"
        - "green"
        - "blue"
        - "yellow"
```

### Single Select

```yaml
questions:
  - name: "License"
    prompt:
      message: "License of the project"
      default: "MIT"
      options:
        - "MIT"
        - "Apache-2.0"
        - "GPL-3.0"
        - "BSD-3-Clause"
        - "Unlicense"
```
