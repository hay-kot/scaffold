---
---

# Scaffold Config Reference

## `questions`

Questions are used to prompt the user for input when generating a scaffold. We support the following types of questions as determined by the `prompt` field.

- [Text Input](#text-input)
- [Multiline Text Input](#multiline-text-input)
- [Looped Text Input](#looped-text-input)
- [Confirm Input](#confirm-input)
- [Select One Input](#select-one-input)
- [Multi Select Input](#multi-select-input)

**Questions have the following properties**

### `name`

The name of the variable that will be used in the template. <span v-pre>`{{ .Scaffold.name }}`</span>

### `required`

Whether or not the question is required.

### `when`

A go template will will be evaluated with the previous context to conditionally render the questions. If the template evaluates to `false` the question will not be rendered, otherwise it will be. This is done by using the `strconv.ParseBool` function to parse the result of the template. 

::: tip
Previous question variables are available at the root level <span v-pre>`{{ .previous_name }}`</span> instead of inside the `.Scaffold` namespace.
::: 

### `group`

The group field is used to group questions together in the rendered form. When inputs share the same group, they are show together in the UI and can be navigated between before submitting the section.

::: warning
When `group` is used in conjunction with the `when` property, only the first question in the group will be evaluated for the `when` field. The result of that evaluation will be applied to the entire group. You cannot apply the `when` condition to individual fields within a group. If you need to filter specific questions, you cannot use groups.
:::

### `prompt`

The prompt field is the type of question to display to the user. See the examples below for more details.

all prompts support the following fields

- `message` - The message to display to the user, think of this as the input label
- `description` - A description of the input, this is displayed below the input
- `default` - The default value for the input, type varies by input type

#### Text Input

Text inputs are the most common and simplest type of inputs, they prompt the user for a text input. The following example will prompt the user for a description of the project.

```yaml
questions:
  - name: "description"
    prompt:
      message: "Description of the project"
```

#### Multiline Text Input

Multiline text inputs allow the user to provide longer text with newlines. The following example will prompt the user for a description of the project.

```yaml
questions:
  - name: "description"
    prompt:
      message: "Description of the project"
      multi: true
```

#### Looped Text Input

Looped text inputs present the same as text inputs, but their resulting type is an array of strings. The following example will prompt the user for a list of CLI commands to stub out.

```yaml
questions:
  - name: "CLI Commands"
    prompt:
      message: "CLI Commands"
      description: "Enter a list of cli commands to stub out"
      loop: true
```

#### Confirm Input

Confirm inputs prompt the user for a yes/no input. The following example will prompt the user to use Github Actions for CI/CD.

```yaml
questions:
  - name: "use_github_actions"
    prompt:
      confirm: "Use Github Actions for CI/CD?"
```

#### Select One Input

```yaml
questions:
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
```

#### Multi Select Input

```yaml
questions:
  - name: "colors"
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

## `computed`

Computed variables are variables that are computed from the answers to the questions. The following example will compute the `shuffled` variable from the `Project` variable.

```yaml
computed:
  shuffled: "{{ shuffle .Project }}"
```

You can reference computed variables like so

```
{{ .Computed.shuffled }}
```

::: tip
Computed variables are generally of type string however, there are special cases for boolean and integer types. Scaffold will attempt to parse the computed string value into an integer, and then a boolean. If any are successful, that value will be used in-place of the string.
:::

## `rewrites`

Rewrites working with the [template scaffolds](../introduction/terminology.md) to perform a path rewrite to another directory. The following example defines a rewrite that will render the `templates/defaults.yaml` file to the <span v-pre>`roles/{{ .ProjectKebab }}/defaults/main.yaml`</span> path.

```yaml
rewrites:
  - from: templates/defaults.yaml
    to: roles/{{ .ProjectKebab }}/defaults/main.yaml
```

- `from` - The path to the template file
- `to` - a template path to the destination file
- These files _are_ rendered with the template engine

This feature is not available for project scaffolds.

## `skips`

Skip is a list of glob patterns that will be used to skip the template **rendering** process. This is useful is your file is itself a go template, or contains similar syntax that will cause the template engine to fail. The following example will skip the `templates/defaults.yaml` file from being rendered.

```yaml
skip:
  - "*.goreleaser.yaml"
  - "**/*.gotmpl"
```

## `inject`

`inject` is a list of code/text injections to perform on a given file. This is to be used in conjunction with `scaffold templates` and is not supported within a `scaffold project`.

The following example will inject a role into the `site.yaml` file at the output directory.

### `name`

The name of the injection

### `path`

The relative path to the file to inject into from the output directory

### `at`

The location to inject the code/text. This is evaluated using the strings.Contains function. Note that ALL matches will be replaced.

### `template`

The template to inject into the file. These work the same as scaffold templates.

**Example**

```yaml
inject:
  - name: "add role to site.yaml"
    path: site.yaml
    at: "# $Scaffold.role_name"
    template: |
      - name: {{ .Scaffold.role_name }}
        role: {{ .Computed.snaked }}
```

## `messages`

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

## `features`

Scaffold support the concept of "feature flags" that can be used to conditionally render entire directories/glob matches of files. This is useful if you want to provide a scaffold that can have wide-reaching optional features, like a database, CI pipeline, etc.

Features are lists of globs that will be evaluated against the `value` field. If the value is parsed to `true` the glob will be rendered, otherwise it will be skipped.

```yaml
features:
  - value: "{{ .Scaffold.database }}"
    globs:
      - "**/core/database/**/*"
```

## `presets`

Presets are a way to define a set of default values for a scaffold. These can be overridden by the user when running the scaffold.

```yaml
presets:
  default:
    description: "A description of the project"
    license: "MIT"
    use_github_actions: true
    colors: ["red", "green"]
```

::: tip Presets and Testing
Presets can be used in conjunction with the `new` command for testing purposes. See [Testing Scaffolds](./testing-scaffolds.md) for more information.
:::
