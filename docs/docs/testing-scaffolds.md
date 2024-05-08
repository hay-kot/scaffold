---
title: Testing Scaffolds
---

The Scaffold CLI has built in support for testing scaffold. The `scaffold test` command will run a case provided with the `--case` flag and render the templates as if those arguments were provided via the interactive interface. If a 'Project' key is not provided in the test case one will be generated in the pattern `scaffold-test-*`.

The `test` command can also output an AST of the rendered scaffold to stdout.

```bash
scaffold test --log-level="panic" --case="default" --memfs --ast <scaffold-name>
scaffold-test-3811:  (type=dir)
        main.go:  (type=file)
                package main

                import (
                        "fmt"
                )

                func main() {
                        fmt.Println("Hello, World!")
                }
```

## Test Cases

Test cases can be defined in the `scaffold.yaml` file under the `tests` key.

```yaml
tests:
  default:
    Var1: "Hello, World!"
    Var2: "Hello, World!"
```

The `default` key is the name of the test case. The values under the test case key are the values that will be passed to the scaffold when running the test case.
