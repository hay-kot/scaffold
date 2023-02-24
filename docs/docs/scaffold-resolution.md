---
title: Scaffold Resolution
---

Scaffold uses a resolution system to find the correct scaffold to use when generating a project. This system may seem _complicated_ at first, however I reason that it is flexible enough to handle most use cases without being overly un-intuitive.

``` mermaid
graph TD;
  A(Invoke Command) --> B(Expand Aliases);
  B --> C{Is Remote URL?}
  C --> |Yes| D{Already Cloned?};
  D --> |Yes| Z;
  D --> |No| E(Clone Repository);
  E --> Z;

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