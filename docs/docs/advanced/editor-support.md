---
---

# Editor Support

Scaffold publishes JSON schema files for both the **scaffoldrc** and **scaffold** config files. This allows for editor to provide intellisense and validation for these files.

To add support for these files in your editor, you can include special comments at the top of the file that point to the schema file.

::: tip
These schema's are a new effort to improve the developer experience for Scaffold users. If you have any feedback or suggestions, please open an issue on the [Scaffold GitHub repository](https://github.com/hay-kot/scaffold/issues/new)
:::

## Scaffold RC

```yaml
# yaml-language-server: $schema=https://hay-kot.github.io/scaffold/scaffoldrc.schema.json
settings:
  theme: scaffold
  run_hooks: prompt
```

## Scaffold Config

```yaml
# yaml-language-server: $schema=https://hay-kot.github.io/scaffold/schema.json
questions: [...]
```
