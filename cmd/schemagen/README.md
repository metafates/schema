# Schemagen

Schemagen is a tool to generate effective implementation of Validateable
interface for your types to reduce validation overhead to **ZERO**.

Whether you choose to use code generation or runtime reflection traversal, validation will work either way.

Using reflection is easier because it does not require any codegen setup, but it does introduce minor performance decrease.

Unless performance is top-priority and validation is indeed a bottleneck (usually it's not), I'd recommend sticking with the reflection - it makes your codebase simpler to maintain.

## How to use

**TODO**: this chapter is not done yet =)

Using `go generate`:

```go
//go:generate schemagen -type MyStruct,MyMap,MySlice

type MyStruct struct {
    Foo required.NonEmpty[string]
    Bar optional.Negative[int]
}

type MyMap map[string]MyStruct

type MySlice map[string]MyStruct
```

```bash
go generate ./...
```

## What does it do

It generates `YOUR_TYPE.schema.gen` file with `Validate() error` method for each type specified.
Therefore `validate.Validate(v any) error` will call this method instead of reflection-based field traversal.

That's it! It will reduce validataion overhead to almost zero.

See [example](./example)
