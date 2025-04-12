# Schemagen

Schemagen is a tool to generate effective implementation of Validateable
interface for your types to reduce validation overhead to **ZERO**.

Whether you choose to use code generation or runtime reflection traversal, validation will work either way.

Using reflection is easier because it does not require any codegen setup, but it does introduce minor performance decrease.

Unless performance is top-priority and validation is indeed a bottleneck (usually it's not), I'd recommend sticking with the reflection - it makes your codebase simpler to maintain. Though I've tried to make this tool as painless to use as go allows =)

## How to use

**WIP**: no stable version yet

<details>
<summary>Go 1.24+ (with tool directive)</summary>

```bash
go get -tool github.com/metafates/schema/cmd/schemagen@main
```

This will add a tool directive to your `go.mod` file

Then you can use it with `go:generate` directive (notice the `go tool` prefix)

```go
//go:generate go tool schemagen -type Foo,Bar

type Foo struct {
    A required.NonEmpty[string]
    B optional.Negative[int]
}

type Bar map[string]MyStruct

type MySlice map[string]MyStruct
```

</details>

<details>
<summary>Go 1.23 and earlier</summary>

See https://marcofranssen.nl/manage-go-tools-via-go-modules

Or:

```bash
go install github.com/metafates/schema/cmd/schemagen@main
```

Ensure that `schemagen` is in your `$PATH`:

```bash
which schemagen # should output something if everything is ok
```

Then you can use it with `go:generate` directive

```go
//go:generate schemagen -type Foo,Bar

type Foo struct {
    A required.NonEmpty[string]
    B optional.Negative[int]
}

type Bar map[string]MyStruct

type MySlice map[string]MyStruct
```

</details>


And call `go generate` as usual

```bash
go generate ./...
```

You should see the following files generated:

- `Foo.schema.go`
- `Bar.schema.go`

## What does it do

It generates `YOUR_TYPE.schema.gen` file with `Validate() error` method for each type specified.
Therefore `validate.Validate(v any) error` will call this method instead of reflection-based field traversal.

That's it! It will reduce validataion overhead to almost zero.
