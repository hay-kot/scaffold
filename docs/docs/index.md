---
title: About
---

!!! warning
    This project is currently unreleased. This documentation is inaccurate and should be viewed as a road map for the project.

Scaffold is a project generation tool similar to [cookiecutter](https://github.com/cookiecutter/cookiecutter) written in Go that leverages the Go template engine to generate projects from a template. It is designed to be used with git or the local file system with a built in "scaffold" management system for versions and templates.

What set's Scaffold apart from projects like cookiecutter is the ability to define reusable scaffolds within a project to help bootstrap code changes in new projects. You're able to use a `.scaffolds` directory within a project to define a scaffold that can generate files in multiple locations around your project. This is useful for generating boilerplate starter code for:

- React/Vue/Angular/... components
- Ansible Roles
- CRUD API endpoint stubs

See the [examples](#examples) section for more information on leveraging this feature.

**Basic Usage**

```sh
scaffold new <scaffold> [flags]
```

See scaffold --help for all available commands and flags


## Featured Scaffolds

### Go CLI

[github.com/hay-kot/scaffold-go-cli](https://github.com/hay-kot/scaffold-go-cli)

- CI/CD with Github Actions or Drone.io
    - PR/Commit/Release workflows
- GoReleaser for releases
- GolangCI-Lint for linting
- Build/Version/Commit injection on build