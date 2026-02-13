# Variable Syntax Reference

Variables are passed as positional arguments to `scaffold new` in the format `key[:type]=value`.

## Syntax

```
key=value            # string (default type)
key:type=value       # explicit type
key:[]type=a,b,c     # slice type (comma-separated)
key:json={"k":"v"}   # arbitrary JSON
```

The first `=` splits key from value. Additional `=` characters in the value are preserved.

## Scalar Types

| Type hint | Aliases | Go type | Example |
|-----------|---------|---------|---------|
| `string` | `str` | `string` | `name=hello` or `name:str=hello` |
| `int` | — | `int` | `port:int=8080` |
| `int32` | — | `int32` | `count:int32=42` |
| `int64` | — | `int64` | `big:int64=999999` |
| `float` | `float64` | `float64` | `rate:float=3.14` |
| `float32` | — | `float32` | `small:float32=1.5` |
| `bool` | — | `bool` | `debug:bool=true` |

When no type is specified, the value is treated as a `string`.

### Boolean values

Accepts: `1`, `t`, `T`, `TRUE`, `true`, `True`, `0`, `f`, `F`, `FALSE`, `false`, `False`

## Slice Types

Slice values are comma-separated. Use `\,` to include a literal comma.

| Type hint | Aliases | Go type | Example |
|-----------|---------|---------|---------|
| `[]string` | `[]str` | `[]string` | `tags:[]string=web,api` |
| `[]int` | — | `[]int` | `ports:[]int=80,443` |
| `[]int32` | — | `[]int32` | — |
| `[]int64` | — | `[]int64` | — |
| `[]float` | `[]float64` | `[]float64` | `rates:[]float=1.1,2.2` |
| `[]float32` | — | `[]float32` | — |
| `[]bool` | — | `[]bool` | `flags:[]bool=true,false` |

An empty value (e.g., `tags:[]string=`) produces an empty slice.

## JSON Type

```bash
config:json='{"db":"postgres","port":8080}'
```

The `json` type parses the value as arbitrary JSON. Supports objects, arrays, strings, numbers, booleans, and null.

## Escaping

- **Commas in slices**: Use `\,` to include a literal comma in a slice element
  ```bash
  names:[]string=Smith\, John,Doe\, Jane
  # Produces: ["Smith, John", "Doe, Jane"]
  ```
- **Backslashes**: Use `\\` for a literal backslash
- A trailing backslash (no character after it) is an error

## Precedence

When using `--preset` with CLI variables:

1. Preset values are loaded first
2. CLI `key=value` arguments are merged on top
3. **CLI arguments win** over preset values
