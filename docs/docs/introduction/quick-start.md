---
---

# Installation

Scaffold is a single-binary written in Go, it can be installed via Homebrew or Go. We also publish binaries as apart of our [GitHub release](https://github.com/hay-kot/scaffold/releases/latest).

### Homebrew

```sh
brew tap hay-kot/scaffold-tap
brew install scaffold
```

### Go

```sh
go install github.com/hay-kot/scaffold@latest
```

## Usage

```sh
scaffold new <scaffold> [flags]
```

Where `<scaffold>` is [a path or a URL to a scaffold](../user-guide/scaffold-resolution.md)

See scaffold --help for all available commands and flags

