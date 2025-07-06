---
---

# Quick Start Guide

## Installation

Scaffold is a single-binary written in Go, it can be installed via Homebrew or Go. We also publish binaries as apart of our [GitHub release](https://github.com/hay-kot/scaffold/releases/latest).

### Homebrew

```sh
brew tap hay-kot/scaffold-tap
brew install scaffold --cask
```

### Go

```sh
go install github.com/hay-kot/scaffold@latest
```

## For Users

Scaffold allows you to quickly generate new projects from templates. You can use scaffolds from GitHub repositories, local files, or built-in templates.

### Basic Usage

```sh
scaffold new <scaffold> [flags]
```

Where `<scaffold>` can be:

- A GitHub repository URL: `scaffold new github.com/username/repo`
- A local path: `scaffold new ./my-local-scaffold`
- A built-in scaffold: `scaffold new hello`

### Using a GitHub Repository

You can directly use scaffolds hosted on GitHub:

```sh
# Create a new project using a GitHub repository
scaffold new github.com/username/repo-name

# Create with a specific version
scaffold new github.com/username/repo-name@v1.0.0

# Create with a specific branch
scaffold new github.com/username/repo-name@branch-name
```

### Listing Available Scaffolds

To see scaffolds you've previously used:

```sh
scaffold list
```

For more detailed information about using scaffolds, see the [Using Projects](../projects/using-projects.md) guide.

## For Creators

Creating your own scaffolds allows you to automate repetitive project setup tasks and share templates with your team or the community.

### Key Features for Scaffold Creators

- **Flexible Structure**: Create unlimited directory nesting with template-based naming
- **Template Engine**: Uses Go template syntax with numerous helper functions
- **Variable Injection**: Define variables that users can input during project creation
- **Partials Support**: Create reusable template components
- **Conditional Features**: Include or exclude files based on user selections
- **Custom Delimiters**: Change template delimiters for specific file types

### Creating Your First Scaffold

1. Create a new directory for your scaffold
2. Add a `scaffold.yml` file to define your configuration
3. Create template files and directories

Basic scaffold structure:

```
my-scaffold/
├── scaffold.yml
├── partials/
│   └── header.txt
└── {{ .Project }}/
    ├── README.md
    └── src/
        └── main.go
```

Basic `scaffold.yml`:

```yaml
questions:
  - name: license
    type: select
    message: Select a license
    options:
      - MIT
      - Apache 2.0
      - GPL
```

For detailed information on creating scaffolds, see:

- [Creating Projects](../projects/creating-projects.md)
- [Scaffold File Configuration](../configuration/scaffold-file.md)
- [Template Engine](../template-system/template-engine.md)

## Additional Commands

See all available commands and options:

```sh
scaffold --help
```
