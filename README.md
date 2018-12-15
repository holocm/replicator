# Replicator

Replicator takes in a template on standard input, reads variables from `/etc/replicator.d` and uses these to render the
template to standard output.

* The template language is that of Go's `text/template` library, so refer to [its
  documentation](https://golang.org/pkg/text/template/) for syntax and available functions. Also, the [Sprig library of
  functions](https://github.com/MasterMinds/sprig) is available. Refer to the release notes for which Sprig version is
  available in a specific release of Replicator.

  The additional template function `toToml` is available, and works analogously to `toJson` from Sprig.

* Variables come from [TOML files](https://github.com/toml-lang/toml) in files
  called `/etc/replicator.d/*.toml`. If there are multiple files, they are
  parsed in alphabetical order and merged as described in the section "Merging"
  below.

For example:

```bash
$ cat /etc/replicator.d/*.toml
[[item]]
name = "foo"

[[item]]
name = "bar"

$ cat example.txt
{{ range .Vars.item }}Hello {{.Name}}. {{end}}

$ replicator < example.txt
Hello foo. Hello bar.
```

## Installation

Replicator requires [Go](https://golang.org) as build-time dependencies. There
are no runtime dependencies other than a libc. Once you're all set, the build
is done with

```
make
sudo make install
```

Replicator can also be installed via `go get github.com/holocm/replicator`.

## Merging

Each TOML file in `/etc/replicator.d/*.toml` is parsed in alphabetical order,
and then they are all merged in the `foldl` pattern:

```
result = (((first + second) + third) + ...)
```

Arrays are merged by concatenation.

```toml
[[first]]
name = "foo"

[[second]]
name = "bar"

[[result]]
name = "foo"
[[result]]
name = "bar"
```

Tables are merged key-by-key, such that values in the second table take precedence
over those with the same key in the first table.

```toml
[first]
key1 = "value1"
key2 = "value2"

[second]
key1 = "value11"
key3 = "value33"

[result]
key1 = "value11"
key2 = "value2"
key3 = "value33"
```

However, if the value from the first table is another table or an array, it
will be merged with the value from the second table by applying these rules
recursively. A type error will be issued if both values have different types
(where all scalar values are considered to be of the same type; so there are
only three types, "scalar", "array" and "table").
