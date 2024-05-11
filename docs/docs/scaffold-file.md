---
title: Scaffold File
---

## Overview

There are two types of scaffolds you can define

### Project Scaffolds

The Project generation scaffolds are used to generate a new project from a template directory. The file structure requires a root directory with one of the following names. All files within that directory will be copied to the destination directory and rendered as a template with the [template engine](./template-engine.md).

- {{ .Project }}
- {{ .ProjectSlug }}
- {{ .ProjectSnake }}
- {{ .ProjectKebab }}
- {{ .ProjectCamel }}

```
├── scaffold # can be any name
    ├── scaffold.yaml
    └── {{ .Project }} # can be any of the project name formats
        └── any nested amount of files...
```

### Template Scaffolds

The template scaffolds are used to generate files within an existing project. The file structure requires a `templates` file in the root directory. The `templates` directory is used to store files that should be rewritten using the [Rewrites](#rewrites) configuration in the `scaffold.yaml` file.

```
├── .scaffolds # in your project directory
    └── my-scaffold # can be any name
        ├── scaffold.yaml
        └── templates
            └── any nested amount of files...
```

The templates directory is _usually_ a flat directory structure, but can be nested as well.

## File Reference

### Questions and Prompts

Questions are used to prompt the user for input when generating a scaffold. We support the following types of questions

- Text
- Confirm (Yes/No)
- Select One
- Multi Select

They share a base type of question with the following fields

`name`

: The name of the variable that will be used in the template. {{ .Scaffold.<name> }}

`required`

: Whether or not the question is required.

`when`

: A go template will will be evaluated with the previous context to conditionally render the questions. If the template evaluates to `false` the question will not be rendered, otherwise it will be. This is done by using the `strconv.ParseBool` function to parse the result of the template. **Previous question variables are available at the root level {{ .previous_name }} instead of inside a .Scaffold container.**

`prompt`

: Prompt configured the type of questions to display to the user. See the examples below for more details.

     `message`

     :   The message field is the message that will be displayed to the user. If only this field is specified, the user will be prompted for a text input.

    `description`

     :   The description field is an optional field that will be displayed to the user as a description of the question.

    `loop`

     :   When the loop field is true, and the question is a text question, the user will be prompted to enter multipel values until they enter an empty string. The resulting type will be an array of string.

     `confirm`

     :   The confirm field is the message that will be displayed to the user. If only this field is specified, the user will be prompted for a yes/no input.

     `options`

     :   The options field is a list of options that will be displayed to the user. This requires the message field to be specified as well.

     `multi`

     :   The multi field is a boolean that will allow the user to select multiple options. This requires the message and options fields to be specified as well.

     `default`

     :    The default field is the default value(s) that will be used if the user does not provide an answer.

#### Question Examples

```yaml
questions:
  - name: "description"
    prompt:
      message: "Description of the project"
    required: true
  - name: "CLI Commands"
    prompt:
      message: "CLI Commands"
      description: "Enter a list of cli commands to stub out"
      loop: true
  - name: "license"
    prompt:
      message: "License of the project"
      default: "MIT"
      options:
        - "MIT"
        - "Apache-2.0"
        - "GPL-3.0"
        - "BSD-3-Clause"
        - "Unlicense"
  - name: "use_github_actions"
    prompt:
      confirm: "Use Github Actions for CI/CD?"
  - name: "colors"
    when: { { .use_github_actions } }
    prompt:
      multi: true
      message: "Colors of the project"
      default: ["red", "green"]
      options:
        - "red"
        - "green"
        - "blue"
        - "yellow"
```

### Computed Variables

Computed variables are variables that are computed from the answers to the questions. The following example will compute the `shuffled` variable from the `Project` variable.

```yaml
computed:
  shuffled: "{{ shuffle .Project }}"
```

You can reference computed variables like so

```yaml
{ { .Computed.shuffled } }
```

!!! tip
    Computed variables are generally of type string however, there are special cases for boolean and integer types. Scaffold will attempt to parse the computed string value into an integer, and then a boolean. If any are successful, that value will be used in-place of the string.



### Rewrites

Rewrites working with the [template scaffolds](./index.md#terminology) to perform a path rewrite to another directory. The following example defines a rewrite that will render the `templates/defaults.yaml` file to the `roles/{{ .ProjectKebab }}/defaults/main.yaml` path.

```yaml
rewrites:
  - from: templates/defaults.yaml
    to: roles/{{ .ProjectKebab }}/defaults/main.yaml
```

- `from` - The path to the template file
- `to` - a template path to the destination file
- These files _are_ rendered with the template engine

This feature is not available for project scaffolds.

### Skips

Skip is a list of glob patterns that will be used to skip the template **rendering** process. This is useful is your file is itself a go template, or contains similar syntax that will cause the template engine to fail. The following example will skip the `templates/defaults.yaml` file from being rendered.

```yaml
skip:
  - "*.goreleaser.yaml"
  - "**/*.gotmpl"
```

### Inject

`inject` is a list of code/text injections to perform on a given file. This is to be used in conjunction with `scaffold templates` and is not supported within a `scaffold project`.

The following example will inject a role into the `site.yaml` file at the output directory.

`name`

: The name of the injection

`path`

: The relative path to the file to inject into from the output directory

`at`

: The location to inject the code/text. This is evaluated using the strings.Contains function. Note that ALL matches will be replaced.

`template`

: The template to inject into the file. These work the same as scaffold templates.

```yaml
inject:
  - name: "add role to site.yaml"
    path: site.yaml
    at: "# $Scaffold.role_name"
    template: |
      - name: {{ .Scaffold.role_name }}
        role: {{ .Computed.snaked }}
```

### Messages

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

### Features

Scaffold support the concept of "feature flags" that can be used to conditionally render entire directories/glob matches of files. This is useful if you want to provide a scaffold that can have wide-reaching optional features, like a database, CI pipeline, etc.

Features are lists of globs that will be evaluated against the `value` field. If the value is parsed to `true` the glob will be rendered, otherwise it will be skipped.

```yaml
features:
  - value: "{{ .Scaffold.database }}"
    globs:
      - "**/core/database/**/*"
```

### Presets

Presets are a way to define a set of default values for a scaffold. These can be overridden by the user when running the scaffold.

```yaml
presets:
  default:
    description: "A description of the project"
    license: "MIT"
    use_github_actions: true
    colors: ["red", "green"]
```

!!! tip "Presets and Testing"
    Presets can be used in conjunction with the `new` command for testing purposes. See [Testing Scaffolds](./testing-scaffolds.md) for more information.
