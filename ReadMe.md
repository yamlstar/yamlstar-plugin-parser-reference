# YAMLStar reference parser plugin

This repository provides the generated Go distribution of the YAML reference
parser maintained in
[`yaml/yaml-reference-parser-clj`](https://github.com/yaml/yaml-reference-parser-clj).

The Clojure project remains the canonical parser source.
This module compiles its released source artifact with Gloat and exposes the
result as an ordinary Go package.
Consumers do not need Gloat, Clojure, a C compiler, or a shared library.

```go
events, err := parser.Parse([]byte("answer: 42\n"))
```

Version `v0.2.5` is generated from `org.yamlstar/yaml-parser` version `0.2.5`.
The source release uses the `v0.2.5` Git tag in
`yaml/yaml-reference-parser-clj`.
This generated module also uses only the `v0.2.5` Git tag.
Go consumers can select either `0.2.5` or `v0.2.5` in plugin configuration,
while documentation and release tags use the `v` prefix.

To regenerate from a local source checkout before that artifact is published:

```sh
make generate YAML_PARSER_DIR=../yaml-reference-parser-clj
```

Run the complete local validation before tagging a release:

```sh
make check-generated YAML_PARSER_DIR=../yaml-reference-parser-clj
make test
```

Release the version recorded in `Makefile` with:

```sh
make release v=0.2.5
```

The command validates tests and generated sources, publishes a clean local
`main` when it is strictly ahead of `origin/main`, and creates the `v0.2.5`
tag and GitHub release.
It never creates an unprefixed version tag.
Use `d=1` to print the release action without publishing it.
