---
title: Scaffold RC
---

Scaffold RC is the runtime configuration file that can be used to define some default values and perform some basic enhancements to the scaffolding process. Your scaffoldrc path is defined by:

- The `--scaffoldrc` flag
- the `SCAFFOLDRC` environment variable

Default: `~/.scaffold/scaffoldrc.yml`

## Defaults

The `defaults` section allows you to set some default values for the scaffolding process. These can be any key/value string pairs

```yaml
defaults:
  name: Joe Bagadonuts
  github_username: joebagadonuts
  email: joebags@donus.gonuts
```

## Aliases

The `aliases` section allows you to define key/value pairs as shortcuts for a scaffold path. This is useful to shorten a reference for a specific scaffold.

```yaml
aliases:
  api: ~/local-scaffolds/api
  component: ~/local-scaffolds/component
```

Then you can use the alias in the `scaffold` command

```bash
scaffold new api
```

## shorts

The `shorts` section allows you to define expandable text snippets. Commonly these would be used to prefix a URL or path.


```yaml
shorts:
    gh: https://github.com/
    gl: https://gitlab.com/
```

Then you can use the alias in the `scaffold` command

```bash
scaffold new gh:joebagadonuts/my-project
```

Which will expand to

```bash
scaffold new https://github.com/joebagadonuts/my-project
```
