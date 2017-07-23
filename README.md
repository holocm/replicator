# Replicator

Replicator takes in a template on standard input, reads variables from `/etc/replicator.d` and uses these to render the
template to standard output.

* The template language is that of Go's `text/template` library, so refer to [its
  documentation](https://golang.org/pkg/text/template/) for syntax and available functions. Also, the [Sprig library of
  functions](https://github.com/MasterMinds/sprig) is available. Refer to the release notes for which Sprig version is
  available in a specific release of Replicator.

* Variables come from [TOML files](https://github.com/toml-lang/toml) in files called `/etc/replicator.d/*.toml`. All
  files are concatenated together (in alphabetical order) before being parsed. These variables are available during
  rendering as `.Vars`.

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

Replicator requires [Go](https://golang.org) and [Perl](https://perl.org) as
build-time dependencies. There are no runtime dependencies other than a libc.
Once you're all set, the build is done with

```
make
sudo make install
```

Replicator can also be installed via `go get github.com/holocm/replicator`.
