---
---

# Scaffold Resolution

Scaffold uses a resolution system to find the correct scaffold to use when generating a project. This system may seem _complicated_ at first, however I reason that it is flexible enough to handle most use cases without being overly un-intuitive.


## Short summary

The scaffold argument can be one of the following

- A scaffold name without slashes in it, e.g. `my-scaffold`. This will be looked up under
  under `./.scaffold/<name>`
- A local absolute path, e.g. `$HOME/scaffolds/my-scaffold`
- A local relative path, e.g. `./local-scaffolds/my-scaffold` or `../shared-scaffolds/my-scaffold`
- A remote repo url, e.g. `https://github.com/hay-kot/scaffold-go-cli`
- A remote repo subdirectory, e.g. `https://github.com/org/repo#subdirectory`

> Note that the resolved path must have a `scaffold.yaml` or `scaffold.yml` [configuration file](../templates/scaffold-file.md) present

## Full resolution diagram

```mermaid
graph TD;
  A(Invoke Command) --> B(Expand Aliases);
  B --> C{Is Remote URL?}
  C --> |Yes| D{Already Cloned?};
  D --> |Yes| repo_subdir_choice;
  D --> |No| E(Clone Repository);
  E --> repo_subdir_choice;

  repo_subdir_choice{#subdir after URL?}

  repo_subdir_choice --> |Yes| use_repo_subdir(Use repo subdirectory)
  repo_subdir_choice --> |No| use_repo_toplevel(Use repo toplevel)

  use_repo_subdir --> Z
  use_repo_toplevel --> Z

  C --> |No| F{Is Absolute Path};
  F --> |Yes| Z;

  F --> |No| G{Contains '/'};
  G --> |Yes| H(Assume Relative Path);
  H --> Z;

  G --> |No| I(Search Scaffold Dirs for Match);
  I --> |Found| Z;
  I --> |Not Found| J(Error);

  Z[Run Scaffold];
```