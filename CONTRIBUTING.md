# Contributing

Thank you for your interest!

## Prerequisite

Install the following programs

- `go` - Of course =) [installation instructions](https://go.dev/doc/install)
- `just` - [Just a command runner](https://github.com/casey/just)
- `python3` - latest stable version (tested with 3.13). No dependencies are needed. This is for code generation.
- `golangci-lint` - [linter and formatter for go](https://golangci-lint.run/welcome/install/)

To contribute follow the following steps:

1. Fork this repository.
2. Make your changes.
3. Run `just` command (without any arguments). Fix errors it emits, if any.
4. Push your changes to your fork.
5. Make PR.
6. You rock!

## Adding new validators

> If you have any questions after reading this section feel free to open an issue - I will be happy to answer.

Validators meta information are stored in [validators.toml](./validators.toml).

This is done so that comment strings and documentation are generated from
a single source of truth to avoid typos and manual work.

After changing this file run:

```sh
just generate
```

After that change [validators/impl.go](./validators/impl.go) file to add `Validate` method for your new validator.

Again, if you have any questions - feel free to open an issue.
