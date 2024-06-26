---
---

# Testing Scaffolds

Scaffold does not have a test command or framework, _however_ it does provide some tools that can be utilized to implement tests for your scaffolds.

## Testing with ASTs

Scaffold provides a way to output an AST of the scaffolded files. This can be used with a diffing tool to compare the ASTs of the scaffolded files with the expected ASTs to ensure that the scaffolded files are correct.

**Command**

```bash
scaffold \
    --log-level="error" \     # set log level to error to avoid noise
    --output-dir=":memory:" \ # render scaffold in memory
    new \
    --preset="default" \      # use scaffold preset
    --no-prompt \             # disable interactive prompts
    --snapshot="stdout" \     # write snapshot to stdout
    <scaffold>
```

**Output**

```bash
scaffold-test-5781:  (type=dir)
        main.go:  (type=file)
                package main

                import (
                        "fmt"
                )

                func main() {
                        fmt.Println("colors=red, green description=This is a test description")
                }
```

## Testing with Outputs

An alternative approach is to test out output scaffolds by running whatever output is generated by the scaffold. For example, we use this approach to test our scaffolds by generating a Go program and running it to ensure that it compiles and runs as expected.

See [cli.test.sh](https://github.com/hay-kot/scaffold/blob/main/tests/cli.test.sh) for an example of how to test scaffolds using this approach.
