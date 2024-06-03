---
---

# Scaffold Runtime Config

Scaffold RC is the runtime configuration file that can be used to define some default values and perform some basic enhancements to the scaffolding process. Your scaffoldrc path is defined by:

- The `--scaffoldrc` flag
- the `SCAFFOLDRC` environment variable

Defaults to `~/.scaffold/scaffoldrc.yml`

## `settings`

The `settings` section allows you to define some global settings for the scaffolding process. These can be any key/value string pairs

**Example**

```yaml
settings:
  theme: scaffold
```

### `theme`

The Theme settings allows the user to set the default theme for the scaffolding process. Options include

- `scaffold` (default)
- `charm`
- `dracula`
- `base16`
- `catppuccino`

### `run_hooks`

You may disable hooks globally by setting `run_hooks` to `never`, or choose to be prompted before they run by setting it to `prompt`. The `--run-hooks` CLI setting takes precedence. Options include:

- `always` - run hooks without prompting
- `never` - never run hooks without prompting
- `prompt` - prompt before running hooks (default)

**Example**

```yaml
settings:
  run_hooks: prompt
```

## `defaults`

The `defaults` section allows you to set some default values for the scaffolding process. These can be any key/value string pairs

```yaml
defaults:
  name: Joe Bagadonuts
  github_username: joebagadonuts
  email: joebags@donus.gonuts
```

## `aliases`

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

## `shorts`

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

## `auth`

The `auth` sections lets you define authentication matchers for your scaffolds. This is useful for using scaffolds that are stored in a private repository.

The configuration supports basic authentication and token authentication. Note that in most cases, you want basic authentication, even us you're using a personal access token.

```yaml
auth:
  - match: github.com/private-repo/*
    basic:
      username: joebagadonuts
      password: ${GITHUB_PASSWORD} # this will be replaced with the environment variable
  - match: gitea.com/private-repo/*
    token: ${GITEA_TOKEN} # this will be replaced with the environment variable
```

::: tip
the `match` key supports regular expressions giving you a lot of flexibility in defining your matchers.
:::
