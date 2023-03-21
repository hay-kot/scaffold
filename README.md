<div>
  <h1 align="center">
    Scaffold
  </h1>
  <p align="center">
    Project Generation Tool
  </p>
  <a href="https://hay-kot.github.io/scaffold/">
    <h2 align="center">
      Read The Docs
    </h2>
  </a>
</div>


Scaffold is a project generation tool similar to [cookiecutter](https://github.com/cookiecutter/cookiecutter) written in Go that leverages the Go template engine to generate projects from a template. It is designed to be used with git or the local file system with a built in "scaffold" management system for versions and templates.

What set's Scaffold apart from projects like cookiecutter is the ability to define reusable scaffolds called `template scaffolds` within a project to help bootstrap code changes in new projects.

You're able to use a `.scaffolds` directory within a project to define a scaffold that can generate files in multiple locations around your project. This is useful for generating boilerplate starter code for:

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
- [x] Feature flag support for blocking/enabling rendering of entire directories/glob matches
