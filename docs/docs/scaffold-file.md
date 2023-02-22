---
title: Scaffold File
---

## Overview

In general there are two _types_ of scaffolds that are supported

- Project Generation
- Bootstraps

The Project generation scaffolds are used to generate a new project from a template. The file structure of this template is

```
├── scaffold # can be any name
    ├── scaffold.yaml
    └── {{ .Project }} # can be any of the project name formats
        └── any nested amount of files...
```

The bootstrap scaffolds are used to generate files within an existing project. The file structure of this template is

```
├── .scaffolds # in your project directory
    └── scaffold # can be any name
        ├── scaffold.yaml
        └── templates
            └── any nested amount of files...
```

The templates directory is _usually_ a flat directory structure, but can be nested as well. Note that the `templates` directory is skipped during the rewrite process and the files are copied to the corresponding [rewrite](#rewrites) paths defined in the configuration file.

## Questions and Prompts

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

## Rewrites

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

## Skips

Skip is a list of glob patterns that will be used to skip the template **rendering** process. This is useful is your file is itself a go template, or contains similar syntax that will cause the template engine to fail. The following example will skip the `templates/defaults.yaml` file from being rendered.

```yaml
skip:
  - "*.goreleaser.yaml"
  - "**/*.gotmpl"
```

## Computed Variables

Computed variables are variables that are computed from the answers to the questions. The following example will compute the `shuffled` variable from the `Project` variable.

```yaml
computed:
  shuffled: "{{ shuffle .Project }}"
```

You can reference computed variables like so

```yaml
{{ .Computed.shuffled }}
```

## Messages

You can specify messages to show the user that are rendered using the [glamour markdown renderer](https://github.com/charmbracelet/glamour/) to show pre and post messages to the user.

```yaml
messages:
  pre: |
    # Pre Message

    This is a pre message that will be shown to the user before the scaffold is generated.

    Template variables are _NOT_ available in this message.

  post: |
    # Post Message

    This is a post message that will be shown to the user after the scaffold is generated.

    You can use variables just as you would in your templates.

    {{ .ProjectKebab }}
```