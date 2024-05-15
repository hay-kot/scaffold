---
---

# Hooks

Hooks are extensionless files that are stored in the `hooks` subdirectory of your scaffold. They allow you to run scripts at specific points during project generation. They are skipped when the scaffold output directory is an in-memory filesystem or when they are explicitely disabled. The [shebang](<https://en.wikipedia.org/wiki/Shebang_(Unix)>) is mandatory and can be set to any interpreter on your system. Template variables are available in the scripts.

Currently, only the `post_scaffold` hook is implemented.

::: tip Working directory
The scripts' working directory is set to the scaffold output directory.
:::

## `post_scaffold`

The `post_scaffold` hook is executed after the files have been rendered on the disk, but before the `post` message is printed. It is typically used to fix the formatting of generated files.
