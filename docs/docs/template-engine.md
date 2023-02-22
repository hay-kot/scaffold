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

## Engine Rules

The template process also uses the following rules for rendering:

1. Empty files are skipped.
2. Files that render to an empty string are skipped.