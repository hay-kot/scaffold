---
---

# Template Engine

Scaffold uses the Go template engine to generate files. The following variables are available to use in your templates at a top level:

- `Project` - The name of the project
- `ProjectKebab` - The kebab case version of the project name
- `ProjectSnake` - The snake case version of the project name
- `ProjectCamel` - The camel case version of the project name
- `ProjectPascal` - The pascal case version of the project name
- `Scaffold` - a map of the scaffold questions and answers
- `Computed` - a map of computed values as defined in the scaffolds configuration

### Template Function

The templates also make available the `sprout` library of functions. See the [sprout documentation](https://docs.atom.codes/sprout) for more information.

We also provide the following functions that help with rendering templates:

#### `wraptmpl`

::: v-pre
Wraps a string in `{{` and `}}` so it can be used as a template. This can also be accomplished by escaping the template syntax. For example, `{{ "{{ .Project }}" }}` will render as `{{ .Project }}`.

    `{{ wraptmpl "docker_dir" }}` -> `{{ "docker_dir" }}`

    vs

    `{{ "{{ docker_dir }}" }}` -> `{{ docker_dir }}`

::: v-pre

## Engine Rules

The template process also uses the following rules for rendering:

1. Empty files are skipped.
2. Template files that are empty after rendering are not included in the generated project.
3. Empty directories not included in the generated project
