---
title: Overview
---

Scaffold is a project generation tool (see [cookiecutter](https://github.com/cookiecutter/cookiecutter)) written in Go that leverages the Go template engine to generate projects from a template. It is designed to be used with git or the local file system with a built in "scaffold" management system for versions and templates.

What set's Scaffold apart from projects like cookiecutter is the ability to define reusable scaffolds within a project to help bootstrap code changes in new projects. You're able to use a `.scaffolds` directory within a project to define a scaffold that can generate files in multiple locations around your project.

See the [examples](#examples) section for more information.

**Usage**

```
scaffold new <scaffold> [flags]
```

See scaffold --help for all available commands and flags
