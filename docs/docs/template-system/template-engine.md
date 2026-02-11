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
- `Each` - available inside [`each`-expanded](../configuration/scaffold-file.md#each) templates, contains `.Each.Item` (the current item string) and `.Each.Index` (the zero-based iteration index)

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

#### `isPlural`

Returns a boolean, `true` if the input is plural, `false` otherwise.

::: v-pre
    `{{ isPlural "apple" }}` -> `false`
    `{{ isPlural "apples" }}` -> `true`
::: v-pre

#### `isSingular`

Returns a boolean, `true` if the input is singular, `false` otherwise.

::: v-pre
    `{{ isSingular "apple" }}` -> `true`
    `{{ isSingular "apples" }}` -> `false`
::: v-pre

#### `toPlural`

Converts a singular word to its plural form.

::: v-pre
    `{{ toPlural "apple" }}` -> `apples`
::: v-pre

#### `toSingular`

Converts a plural word to its singular form.

::: v-pre
    `{{ toSingular "apples" }}` -> `apple`
::: v-pre

## Engine Rules

The template process also uses the following rules for rendering:

1. Empty files are skipped.
2. Template files that are empty after rendering are not included in the generated project.
3. Empty directories not included in the generated project
