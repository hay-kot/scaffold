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

!!! tip
    Note that defaults are only used for text type questions at this time.

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

## Shorts

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

## Auth

The `auth` sections lets you define authentication matchers for your scaffolds. This is useful for using scaffolds that are stored in a private repository.

The configuration supports basic authentication and token authentication. Note that in most cases, you want basic authnetication, even us you're using a personal access token.

```yaml
auth:
  - match: github.com/private-repo/*
    basic:
      username: joebagadonuts
      password: ${GITHUB_PASSWORD} # this will be replaced with the environment variable
  - match: gitea.com/private-repo/*
    token: ${GITEA_TOKEN} # this will be replaced with the environment variable
```

!!! tip
    the `match` key supports regular expressions giving you a lot of flexibility in defining your matchers.
