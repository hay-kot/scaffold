---
title: Testing Scaffolds
---

Scaffold has some built-in support for testing scaffolds. This is done through with the `--test` flag when running the `scaffold new --test <template>` command, this will bypass the Q/A step of the scaffold and instead use the `test` property from the scaffold file.

## Scaffold Test Property

The `test` property is a map of key-value pairs that will be used to fill in the values of the scaffold. The key is the name of the variable in the scaffold file and the value is the value that will be used to fill in the variable.

```yaml
test:
  Var1: "Value1"
  Val:
    - "Value1"
    - "Value2"
```
