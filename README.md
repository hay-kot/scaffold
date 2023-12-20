<div>
  <img src="/docs/docs/assets/imgs/scaffold-gopher.webp">
  <h1 align="center" style="margin-bottom: 0px;">
    Scaffold
  </h1>
  <p align="center" style="margin-top: -10px;">
    A Project Generation Tool
  </p>
  <div align="center">
    <a href="https://hay-kot.github.io/scaffold/">
      Read The Docs
    </a>
  </div>
</div>

<br />

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


## Credits

- Logo By [@lakotelman](https://github.com/lakotelman)