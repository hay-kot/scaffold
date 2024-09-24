---
---

# Scaffold File

There are two types of scaffolds you can define:

- [Project Scaffolds](#project-scaffolds), to generate new project directories
- [Template Scaffolds](#template-scaffolds), to generate new files in existing projects

Both require that you have a `scaffold.yaml` or `scaffold.yml` config file placed within the scaffold directory.

## Project Scaffolds

The Project generation scaffolds are used to generate a new project from a template directory. The file structure requires a root directory with one of the following names. All files within that directory will be copied to the destination directory and rendered as a template with the [template engine](./template-engine.md).


::: v-pre
- `{{ .Project }}`
- `{{ .ProjectSlug }}`
- `{{ .ProjectSnake }}`
- `{{ .ProjectKebab }}`
- `{{ .ProjectCamel }}`
- `{{ .ProjectPascal }}`
:::

### File Structure

```
├── scaffold # can be any name
    ├── scaffold.yaml
    └── {{ .Project }} # can be any of the project name formats
        └── any nested amount of files...
```

## Template Scaffolds

The template scaffolds are used to generate files within an existing project. The file structure requires a `templates` folder in the root directory. The `templates` directory is used to store files that should be rewritten using the [rewrites](./config-reference#rewrites) configuration in the `scaffold.yaml` file.


### File Structure

```
 my-scaffold # can be any name
 ├── scaffold.yaml
 └── templates
     └── any nested amount of files...
```

The templates directory is _usually_ a flat directory structure, but can be nested as well.


### Project-specific scaffolds

Project-specific template scaffolds can be nested in the project itself:

```
├── .scaffolds # in your project directory
    ├── my-scaffold
        ├── scaffold.yaml
        └── templates
            └── any nested amount of files...
    └── other-scaffold
        ├── scaffold.yaml
        └── templates
            └── any nested amount of files...
```

in which case, they're accessible directly by name e.g. `scaffold new my-scaffold`

However, they can also be placed in a common directory or remote repository and
used by specifying a local path or a remote URL
