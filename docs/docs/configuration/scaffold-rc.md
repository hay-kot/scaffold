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

## `snippets`

The `snippets` section lists where scaffold looks for [snippets](/snippets/usage). Each entry is a directory, a single file, or a glob pattern. Entries support `~` and environment variables, and relative entries are resolved against the directory that holds your scaffoldrc file.

```yaml
snippets:
  - ~/dotfiles/snippets
  - ~/work/templates/**/*.snippet.md
```

When `snippets` is not set, scaffold uses the `snippets` directory next to your scaffoldrc file. If two locations have a snippet with the same name, the one listed first wins.

## `scaffolds`

The `scaffolds` section lists remote scaffolds that are always available, so you don't have to type the remote each time. Each entry uses the same syntax as the `scaffold new` argument: a URL, an SSH URL, or a short, with an optional `@ref` and `#subdir`.

```yaml
shorts:
  gh: https://github.com/

scaffolds:
  - gh:hay-kot/scaffold-go-cli
  - https://github.com/org/templates@v1.2.0#service
```

Before `scaffold new` shows its picker, and when you run `scaffold list` or `scaffold update`, scaffold clones every entry that isn't in the cache yet. Cloned entries show up as system scaffolds. Nothing is pulled at that point -- `scaffold update` pulls the latest changes.

- Clones use the matching [`auth`](#auth) entry. No interactive auth prompt is shown.
- If a clone fails, for example when you are offline, scaffold prints a warning to stderr and carries on with the scaffolds already cached. Press Ctrl+C during the clone spinner to skip a clone that hangs; scaffold stops the clone, removes what it had written, and continues. The entry is cloned on a later run.
- Entries pinned with `@ref` are cloned at that ref and `scaffold update` doesn't pull them.

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
