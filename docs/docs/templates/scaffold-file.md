---
---

# Scaffold File

There are two types of scaffolds you can define.

## Project Scaffolds

The Project generation scaffolds are used to generate a new project from a template directory. The file structure requires a root directory with one of the following names. All files within that directory will be copied to the destination directory and rendered as a template with the [template engine](./template-engine.md).


::: v-pre
- `{{ .Project }}`
- `{{ .ProjectSlug }}`
- `{{ .ProjectSnake }}`
- `{{ .ProjectKebab }}`
- `{{ .ProjectCamel }}`
:::

### File Structure

```
├── scaffold # can be any name
    ├── scaffold.yaml
    └── {{ .Project }} # can be any of the project name formats
        └── any nested amount of files...
```

## Template Scaffolds

The template scaffolds are used to generate files within an existing project. The file structure requires a `templates` file in the root directory. The `templates` directory is used to store files that should be rewritten using the [Rewrites](#rewrites) configuration in the `scaffold.yaml` file.

### File Structure

```
├── .scaffolds # in your project directory
    └── my-scaffold # can be any name
        ├── scaffold.yaml
        └── templates
            └── any nested amount of files...
```

The templates directory is _usually_ a flat directory structure, but can be nested as well.