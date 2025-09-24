---
---

# What is Scaffold?

Scaffold is a project generation tool similar to [cookiecutter](https://github.com/cookiecutter/cookiecutter) written in Go that leverages the Go template engine to generate projects from a template. It is designed to be used with git or the local file system with a built in "scaffold" management system for versions and templates.

What set's Scaffold apart from projects like cookiecutter is the ability to define reusable scaffolds called `template scaffolds` within a project to help bootstrap code changes in existing projects.

## Use Cases

- **Whole Project Scaffolding**

  Scaffold is great for generating whole projects from a template. This is useful for:

  - Bootstrapping a new project
  - Creating a new API using a standard layout
  - Create a new microservice using company standards

- **Templates within Projects**

  You're able to use a `.scaffolds` directory within a project to define a scaffold that can generate files in multiple locations around your project and even inject code into existing files. This is useful for generating boilerplate starter code for:

    - React/Vue/Angular/... components
    - Ansible Roles
    - CRUD API endpoint stubs
    - Other commonly structured code folders

- **Shared remote templates**

  Templates that add files to a project don't have to be nested under `.scaffolds`. For example, if you are building a tool which users can add to existing project and that tool needs configuration, you can host those tool's scaffolds in a remote repository

See the [examples](https://github.com/hay-kot/scaffold/tree/main/.examples) folder for some examples of how to use Scaffold.

## Features

- Generate projects from a template
- Git based scaffolds (public and private)
    - Update scaffolds with `scaffold update`
    - List scaffolds with `scaffold list`
- Generate files in multiple locations within an existing project
- Pre/Post Messages defined in the scaffold (supports markdown)
- Alias support for shortening common commands
- Shortcuts for common prefixes (e.g `gh:` for github.com)
- Conditional Prompting based on user input
- Inject snippets into existing files with Scaffold Templates
