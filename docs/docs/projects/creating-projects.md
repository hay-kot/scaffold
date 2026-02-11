# Creating Scaffolds

This guide explains how to create project scaffolds that others can use to generate new projects.

## Project Scaffold Structure

A project scaffold requires a specific directory structure:

:::v-pre
```
scaffold-directory/           # Can be any name
├── scaffold.yaml             # Configuration file
└── {{ .Project }}/           # Root directory with dynamic name
    ├── file1.txt             # Template files
    ├── file2.md              # Template files
    └── src/                  # Nested directories
        └── more-files...     # More template files
```
:::

## Root Directory Naming

The root directory within your scaffold must use one of the following dynamic naming patterns:

:::v-pre
- `{{ .Project }}` - Base project name
- `{{ .ProjectSlug }}` - Project name in slug format
- `{{ .ProjectSnake }}` - Project name in snake_case
- `{{ .ProjectKebab }}` - Project name in kebab-case
- `{{ .ProjectCamel }}` - Project name in camelCase
- `{{ .ProjectPascal }}` - Project name in PascalCase
:::

## Configuration File

Project scaffolds require a `scaffold.yaml` or `scaffold.yml` config file placed at the root of the scaffold directory. This file has a whole range of options for defining questions to ask the user and different flags/options that can be set.

See [Scaffold Config](../configuration/scaffold-file.md) for all available options.


## Multi-File Output

Scaffold supports generating multiple files or directories from a single template using the `[varname]` path convention. When a directory or filename contains `[varname]` and the variable is declared in the [`each`](../configuration/scaffold-file.md#each) config, the path is expanded once per item in the list.

:::v-pre
```
scaffold-directory/
├── scaffold.yaml
└── {{ .Project }}/
    ├── [services]/          # One directory per service
    │   ├── handler.go
    │   └── routes.go
    └── main.go
```
:::

See the [`each` configuration reference](../configuration/scaffold-file.md#each) for details on setup and the `as` transformation option.

## Template Files

All files within the root directory will be processed by the template engine. You can use template syntax to customize content:

```markdown
# {{ .Project }}

{{ .Scaffold.description }}

## License

This project is licensed under the {{ .Scaffold.license }} License.
```

For more details on template syntax, see the [template engine documentation](../template-system/template-engine.md).