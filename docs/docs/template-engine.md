---
title: Template Engine
---

## Overview

Scaffold uses the Go template engine to generate files. The following variables are available to use in your templates at a top level:

- `Project` - The name of the project
- `ProjectKebab` - The kebab case version of the project name
- `ProjectSnake` - The snake case version of the project name
- `ProjectCamel` - The camel case version of the project name
- `Scaffold` - a map of the scaffold questions and answers

The templates also make available the `sprig` library of functions. See the [sprig documentation](http://masterminds.github.io/sprig/) for more information.

!!! tip "Searching Template Functions"

    At the time of writing this the sprig documentation lacks a search function. I've submitted a [PR #363](https://github.com/Masterminds/sprig/pull/363) to remedy this. In the meantime you can use my forks documentation to search for functions.

    - [https://hay-kot.github.io/sprig/](https://hay-kot.github.io/sprig/)

We also provide the following functions that help with rendering templates:

`wraptmpl`

:    Wraps a string in `{{` and `}}` so it can be used as a template. This can also be accomplished by escaping the template syntax. For example, `{{ "{{ .Project }}" }}` will render as `{{ .Project }}`.

    `{{ wraptmpl "docker_dir" }}` -> `{{ "docker_dir" }}`

    vs

    `{{ "{{ docker_dir }}" }}` -> `{{ docker_dir }}`

## Engine Rules

The template process also uses the following rules for rendering:

1. Empty files are skipped.
2. Template files that are empty after rendering are not included in the generated project.
3. Empty directories not included in the generated project