# Schema Files

In this directory there are two schema files.

1. schema.scaffold.ts
2. schema.scaffoldrc.ts

They define the type types of the configuration files that are used by scaffold.
We write these schemas in typescript and they are compiled into json schema files during the build process for the documentation.
We utilize typescript because we **1)** already have a typescript dependency in the project and **2)** I (hay-kot) much prefer writing the schema in typescript over YAML or raw JSON.

## How to Update Schema Files

When you want to update a configuration property in the `scaffoldrc` or the `scaffold` file, there are multiple steps to this process.

1. Update the Go code to support those new fields
2. Update the Typescript Schema files to support those new fields
3. Update the documentation to reflect the new fields. _This is not automatic and must be done manually_.

Once you've made those updates, they will automatically be included in the bundle/build for the docs and new changes will be available for use via the LSP server for the yaml files.