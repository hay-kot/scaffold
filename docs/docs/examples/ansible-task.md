---
title: Ansible Task
---

Given the following starting structure, we can use the command `scaffold new role` to generate a new role in the `ansible/roles` directory.

**Starting File System**

```sh
├── .scaffolds
│   └── role
│       ├── scaffold.yaml
│       └── templates
│           ├── task.yaml
│           └── defaults.yaml
├── ansible
    └── roles
        └── my-existing
```

**Scaffold File**

```yaml
questions:
  - name: "do"
    prompt:
      confirm: "Ok?"
rewrites:
  - from: templates/defaults.yaml
    to: roles/{{ .ProjectKebab }}/defaults/main.yaml
  - from: templates/task.yaml
    to: roles/{{ .ProjectKebab }}/tasks/main.yaml
```

**Output Directory**

```sh
├── .scaffolds
│   └── role
│       ├── scaffold.yaml
│       └── templates
│           ├── task.yaml
│           └── defaults.yaml
└─── ansible
    └── roles
        ├── my-existing
        └── my-new-role
            ├── defaults
            │   └── main.yaml
            └── tasks
                └── main.yaml
```