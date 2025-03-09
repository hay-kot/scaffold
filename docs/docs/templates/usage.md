---
---

# Using Template Scaffolds

This guide explains how to use template scaffolds to generate files within an existing project using templates in your project's `.scaffolds` directory.

## Basic Usage

To generate new files using a template scaffold, use the `scaffold new` command followed by the template name:

```bash
# Generate files using a template
scaffold new component
```

## Interactive Prompts

When you run the command, you'll enter an interactive prompt that asks questions required to render the template:

```
? Component Name: UserProfile
? Component Type: Functional
? Include Tests: Yes
```

The answers you provide will customize the generated files according to the template's design.

## Output Location

By default, template scaffolds will generate files in your current directory according to the rules defined in the template.

You can override the base output directory using the `--output-dir` flag:

```bash
scaffold new --output-dir ./src/components component
```

## Project-Specific Templates

Project-specific templates are stored in the `.scaffolds` directory of your project:

```
project/
├── .scaffolds/
│   ├── component/
│   │   ├── scaffold.yaml
│   │   └── templates/
│   │       └── ...
│   └── page/
│       ├── scaffold.yaml
│       └── templates/
│           └── ...
└── src/
    └── ...
```

To use these templates, simply provide the template name:

```bash
scaffold new component
```

*A full list of flags and options is available with* `scaffold new --help`