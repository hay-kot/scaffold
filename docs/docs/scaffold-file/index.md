---
title: Overview
---

In general there are two _types_ of scaffolds that are supported

- Project Generation
- Bootstraps

The Project generation scaffolds are used to generate a new project from a template. The file structure of this template is

```
├── scaffold # can be any name
    ├── scaffold.yaml
    └── {{ .Project }} # can be any of the project name formats
        └── any nested amount of files...
```

The bootstrap scaffolds are used to generate files within an existing project. The file structure of this template is

```
├── .scaffolds # in your project directory
    └── scaffold # can be any name
        ├── scaffold.yaml
        └── templates
            └── any nested amount of files...
```

The templates directory is _usually_ a flat directory structure, but can be nested as well. Note that the `templates` directory is skipped during the rewrite process and the files are copied to the corresponding [rewrite](#rewrites) paths defined in the configuration file