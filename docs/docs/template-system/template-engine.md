---
---

# Template Engine

Scaffold uses the Go template engine to generate files. The following variables are available to use in your templates at a top level:

- `Project` - The name of the project
- `ProjectKebab` - The kebab case version of the project name
- `ProjectSnake` - The snake case version of the project name
- `ProjectCamel` - The camel case version of the project name
- `ProjectPascal` - The pascal case version of the project name
- `Scaffold` - a map of the scaffold questions and answers
- `Computed` - a map of computed values as defined in the scaffolds configuration
- `Each` - available inside [`each`-expanded](../configuration/scaffold-file.md#each) templates, contains `.Each.Item` (the current item string) and `.Each.Index` (the zero-based iteration index)
- `Ctx` - facts about where scaffold is running, see [Render Context](#render-context)

## Render Context

`Ctx` carries the filesystem and repository facts that differ between checkouts.
Use it when a template has to produce a different value in each copy of a project,
such as a port block or a docker stack name.

- `Ctx.OutputDir` - the output directory as given on the command line
- `Ctx.OutputDirAbs` - the absolute output directory
- `Ctx.OutputDirBase` - the last path segment of the output directory
- `Ctx.WorkingDir` - the directory scaffold was started from
- `Ctx.WorkingDirBase` - the last path segment of the working directory
- `Ctx.Git.Branch` - the current branch, empty on a detached HEAD
- `Ctx.Git.Repo` - the directory name of the git worktree root

Every field is best effort. A missing git repository or an unreadable directory
gives an empty string instead of an error, so guard on the value if a template
cannot work without it.

::: tip
`Ctx.OutputDirAbs` and `Ctx.OutputDirBase` are empty when rendering to
`:memory:` or during `--dry-run`, because there is no directory on disk.
:::

### Template Function

The templates also make available the `sprout` library of functions. See the [sprout documentation](https://docs.atom.codes/sprout) for more information.

We also provide the following functions that help with rendering templates:

#### `wraptmpl`

::: v-pre
Wraps a string in `{{` and `}}` so it can be used as a template. This can also be accomplished by escaping the template syntax. For example, `{{ "{{ .Project }}" }}` will render as `{{ .Project }}`.

    `{{ wraptmpl "docker_dir" }}` -> `{{ "docker_dir" }}`

    vs

    `{{ "{{ docker_dir }}" }}` -> `{{ docker_dir }}`

::: v-pre

#### `portBlock`

Maps a seed string to a deterministic, block-aligned base port.

::: v-pre
    `{{ portBlock "recipinned" }}` -> `26112`
    `{{ portBlock "recipinned" 40000 49000 10 }}` -> `45740`
::: v-pre

It takes the seed and optionally `start`, `end`, and `size`. The result is the
first port of a block of `size` consecutive ports, so add an offset per service:

```
{{ $base := portBlock .Ctx.OutputDirBase }}
API_PORT={{ add $base 1 }}
WEB_PORT={{ add $base 2 }}
```

The defaults are `20000`, `32767`, and `16`. The window stops below 32768 because
that is where Linux starts handing out ephemeral source ports; macOS starts at
49152. A port above that line can be claimed by an unrelated outbound connection,
which shows up as an intermittent bind failure.

Align to a block rather than hashing straight to a port and incrementing. Two raw
hashes can land a few ports apart and partially overlap, so some services bind and
others do not. With alignment, two seeds either share the whole block or share
nothing.

::: warning Stable, not unique
The same seed always gives the same block, with no state to store, but distinct
seeds can collide. Collision odds follow the birthday bound: the defaults give
798 blocks, so roughly 5.5% at 10 concurrent seeds and 21% at 20. Widen the range
or shrink the block size if you run many stacks at once. Treat the result as a
good default a human can override, not as a reservation. Writing the base to
a file that [`no_clobber`](../configuration/scaffold-file.md) protects makes the
override stick across re-runs.
:::

#### `isPlural`

Returns a boolean, `true` if the input is plural, `false` otherwise.

::: v-pre
    `{{ isPlural "apple" }}` -> `false`
    `{{ isPlural "apples" }}` -> `true`
::: v-pre

#### `isSingular`

Returns a boolean, `true` if the input is singular, `false` otherwise.

::: v-pre
    `{{ isSingular "apple" }}` -> `true`
    `{{ isSingular "apples" }}` -> `false`
::: v-pre

#### `toPlural`

Converts a singular word to its plural form.

::: v-pre
    `{{ toPlural "apple" }}` -> `apples`
::: v-pre

#### `toSingular`

Converts a plural word to its singular form.

::: v-pre
    `{{ toSingular "apples" }}` -> `apple`
::: v-pre

## Engine Rules

The template process also uses the following rules for rendering:

1. Empty files are skipped.
2. Template files that are empty after rendering are not included in the generated project.
3. Empty directories not included in the generated project
