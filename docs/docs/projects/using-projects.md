# Using Scaffolds

This guide explains how to use existing project scaffolds to generate new projects.

## Basic Usage

To generate a new project using a scaffold, use the `scaffold new` command followed by either a URL or filepath to the scaffold:

```bash
# Generate from a remote repository
scaffold new https://github.com/hay-kot/scaffold-go-cli

# Generate from a local scaffold
scaffold new path/to/local/scaffold
```

## Interactive Prompts

When you run the `scaffold new` command, you'll enter an interactive prompt that will ask you questions required to render the template correctly. These questions are defined in the scaffold's configuration.

```
? Project Name: my-awesome-project
? Description: A description of my awesome project
? Author: Your Name
```

The answers you provide will be used to customize the generated project according to the template's design.

## Output Location

By default, the scaffold will generate files in your current directory. You can specify a different output directory using the `--output-dir` flag:

```bash
scaffold new --output-dir ./my-new-project https://github.com/hay-kot/scaffold-go-cli
```

*A full list of flags and options is available in the CLI with* `scaffold new --help`