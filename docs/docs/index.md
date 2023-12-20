---
title: Scaffold
---

<div style="display: grid; place-content: center;">
    <img src="./assets/imgs/scaffold-gopher.webp" align="center" style="max-height: 350px;" />
</div>


Scaffold is a project generation tool similar to [cookiecutter](https://github.com/cookiecutter/cookiecutter) written in Go that leverages the Go template engine to generate projects from a template. It is designed to be used with git or the local file system with a built in "scaffold" management system for versions and templates.

What set's Scaffold apart from projects like cookiecutter is the ability to define reusable scaffolds called `template scaffolds` within a project to help bootstrap code changes in existing projects.

You're able to use a `.scaffolds` directory within a project to define a scaffold that can generate files in multiple locations around your project and even inject code into existing files. This is useful for generating boilerplate starter code for:

- React/Vue/Angular/... components
- Ansible Roles
- CRUD API endpoint stubs
- Other commonly structured code folders

See the [examples](https://github.com/hay-kot/scaffold/tree/main/.examples) folder for some examples of how to use Scaffold.

## Core Features

- [x] Generate projects from a template
- [ ] Git based scaffolds
    - [x] Update scaffolds with `scaffold update`
    - [x] List scaffolds with `scaffold list`
    - [ ] Pull specific tag or branch
- [x] Generate files in multiple locations within an existing project
- [x] Pre/Post Messages defined in the scaffold (supports markdown)
- [x] Alias support for shortening common commands
- [x] Shortcuts for common prefixes (e.g `gh:` for github.com)
- [x] Conditional Prompting based on user input
- [x] Inject snippets into existing files with Scaffold Templates

## Installation

### Homebrew

```sh
brew tap hay-kot/scaffold-tap

brew install scaffold
```

### Go

```sh
go install github.com/hay-kot/scaffold@latest
```

## Usage

```sh
scaffold new <scaffold> [flags]
```

See scaffold --help for all available commands and flags


## Terminology

Some of the terms used in the documentation and project can be somewhat general, these definitions help clarify the meaning of the terms used.

`scaffold`

:   a generic term for a repository or directory that has a `scaffold.{yml,yaml}` file in it.

`project`

:   a `scaffold` type that is used to generate a _new_ project, it contains one of the special scaffold project directories

`template`

:   a `scaffold` type that uses the rewrite feature to generate files into multiple places. This is used in an existing directory to add new files. You would use a `template scaffold` to generate the boilerplate files for a new Vue component or Ansible role.


## Featured Scaffolds

### Go CLI

[github.com/hay-kot/scaffold-go-cli](https://github.com/hay-kot/scaffold-go-cli)

- CI/CD with Github Actions or Drone.io
    - PR/Commit/Release workflows
- GoReleaser for releases
- GolangCI-Lint for linting
- Build/Version/Commit injection on build