# 📐 Schema

[![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/metafates/schema)
[![Go Report Card](https://goreportcard.com/badge/github.com/metafates/schema)](https://goreportcard.com/report/github.com/metafates/schema)

> Work in progress!

<img src="https://github.com/user-attachments/assets/54aafdf2-df4d-4b92-9a46-37bc59d99e6e" align="right" width=300 />

Go schema declaration and validation with static types.
No field tags or code duplication.

Schema is designed to be as developer-friendly as possible.
The goal is to eliminate duplicative type declarations.
You declare a schema once and it will be used as both schema and type itself.
It's easy to compose simpler types into complex data structures.

No stable version yet, but you can use it like that.

```bash
go get github.com/metafates/schema@main
```

**Work in progress, API may change significantly without further notice!**

## Features

- Type-safe
- Zero setup required
- Zero overhead (can be achieved with optional codegen)
- No DSL or code duplication
- Cross-field validation support
- Helpful errors
- Parse arbitrary types into schema. E.g. validate **gRPC** generated messages by parsing into schema structs.

## Example

See [examples](./examples) for more examples

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	schemajson "github.com/metafates/schema/encoding/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/charset"
)

// Let's assume we have a request which accepts an user
type User struct {
	// User name is required and must not be empty string
	Name required.NonZero[string] `json:"name"`

	// Birth date is optional, which means it could be null.
	// However, if passed, it must be an any valid [time.Time]
	Birth optional.Any[time.Time] `json:"birth"`

	// Same for email. It is optional, therefore it could be null.
	// But, if passed, not only it must be a valid string, but also a valid RFC 5322 email string
	Email optional.Email[string] `json:"email"`

	// Bio is just a regular string. It may be empty, may be not.
	// No further logic is attached to it.
	Bio string `json:"bio"`
}

// We could also have an address
type Address struct {
	// Latitude is required and must be a valid latitude (range [-90; 90])
	Latitude required.Latitude[float64] `json:"latitude"`

	// Longitude is also required and must be a valid longitude (range [-180; 180])
	Longitude required.Longitude[float64] `json:"longitude"`
}

// But wait, what about custom types?
// We might want (for some reason) a field which accepts only short strings (<10 bytes).
// Let's see how we might implement it.

// let's define a custom validator for short strings.
// it should not contain any fields in itself, they won't be initialized or used in any way.
type ShortStr struct{}

// this function implements a special validator interface
func (ShortStr) Validate(v string) error {
	if len(v) >= 10 {
		return errors.New("string is too long")
	}

	return nil
}

// that's it, basically. now we can use this validator in our request.

// but we can go extreme! we can combine multiple validators using types
type ASCIIShortStr struct {
	// both ASCII and ShortStr must be satisfied.
	// you can also use [validate.Or] to ensure that at least one condition is satisfied.
	validate.And[
		string,
		validate.Charset0[string, charset.ASCII],
		ShortStr,
	]
}

// Now, our final request may look something like that
type Request struct {
	// User is required by default.
	// Because, if not passed, it will be empty.
	// Therefore required fields in user will also be empty which will result a missing fields error.
	User User `json:"user"`

	// Address is, however, optional. It could be null.
	// But, if passed, it must be a valid address with respect to its fields validation (required lat/lon)
	Address optional.Any[Address] `json:"address"`

	// this is how we can use our validator in custom type.
	// we could make an alias for that custom required type, if needed.
	MyShortString required.Custom[string, ShortStr] `json:"myShortString"`

	// same as for [Request.MyShortString] but using an optional instead.
	ASCIIShortString optional.Custom[string, ASCIIShortStr] `json:"asciiShortString"`

	PermitBio bool `json:"permitBio"`
}

// - "How do I do cross-field validation?"
// - Implement [validate.Validateable] interface for your struct
//
// This method will be called AFTER required and optional fields are validated.
// It is optional - you may skip defining it if you don't need to.
func (r *Request) Validate() error {
	if !r.PermitBio {
		if r.User.Bio != "" {
			return errors.New("bio is not permitted")
		}
	}

	return nil
}

func main() {
	// here we create a VALID json for our request. we will try to pass an invalid one later.
	data := []byte(`
{
	"user": {
		"name": "john",
		"email": "john@example.com (comment)",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "foo",
	"permitBio": true
}`)

	var request Request

	{
		// let's unmarshal it. we can use anything we want, not only json
		if err := json.Unmarshal(data, &request); err != nil {
			log.Fatalln(err)
		}

		// but validation won't happen just yet. we need to invoke it manually
		// (passing pointer to validate is required to maintain certain guarantees later)
		if err := validate.Validate(&request); err != nil {
			log.Fatalln(err)
		}
		// that's it, our struct was validated successfully!
		// no errors yet, but we will get there
		//
		// remember our custom cross-field validation?
		// it was called as part of this function
	}

	{
		// we could also use a helper function that does *exactly that* for us.
		if err := schemajson.Unmarshal(data, &request); err != nil {
			log.Fatalln(err)
		}
	}

	// now that we have successfully unmarshalled our json, we can use request fields.
	// to access values of our schema-guarded fields we can use Get() method
	//
	// NOTE: calling this method BEFORE we have
	// validated our request will panic intentionally.
	fmt.Println(request.User.Name.Get()) // output: john

	// optional values return a tuple: a value and a boolean stating its presence
	email, ok := request.User.Email.Get()
	fmt.Println(email, ok) // output: john@example.com (comment) true

	// birth is missing so "ok" will be false
	birth, ok := request.User.Birth.Get()
	fmt.Println(birth, ok) // output: 0001-01-01 00:00:00 +0000 UTC false

	// let's try to pass an INVALID jsons.
	invalidEmail := []byte(`
{
	"user": {
		"name": "john",
		"email": "john@@@example.com",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "foo",
	"permitBio": true
}`)

	invalidShortStr := []byte(`
{
	"user": {
		"name": "john",
		"email": "john@example.com",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "super long string that shall not pass!!!!!!!!",
	"permitBio": true
}`)

	missingUserName := []byte(`
{
	"user": {
		"email": "john@example.com",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "foo",
	"permitBio": true
}`)

	bioNotPermitted := []byte(`
{
	"user": {
		"name": "john",
		"email": "john@example.com",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "foo",
	"permitBio": false
}`)

	fmt.Println(schemajson.Unmarshal(invalidEmail, new(Request)))
	// validate: .User.Email: mail: missing '@' or angle-addr

	fmt.Println(schemajson.Unmarshal(invalidShortStr, new(Request)))
	// validate: .MyShortString: string is too long

	fmt.Println(schemajson.Unmarshal(missingUserName, new(Request)))
	// validate: .User.Name: missing value

	fmt.Println(schemajson.Unmarshal(bioNotPermitted, new(Request)))
	// validate: bio is not permitted

	// You can check if it was validation error or any other json error.
	err := schemajson.Unmarshal(missingUserName, new(Request))

	var validationErr validate.ValidationError
	if errors.As(err, &validationErr) {
		fmt.Println("error while validating", validationErr.Path())
		// error while validating .User.Name

		fmt.Println(errors.Is(err, required.ErrMissingValue))
		// true
	}

	// one more feature - parsing!
	// not all data comes from json - sometimes we already have some initialized values (structs, maps) as go values.
	//
	// for example, we may have some generated code with gRPC and we want to validate it.
	// what we can do with this library:

	// first, let's define some structure
	type Example struct {
		ID           required.UUID[string]
		Content      string
		Tags         required.NonEmptySlice[string]
		RandomNumber float64
	}

	// second, let's assume we have the following grpc message generated
	type GRPCExample struct {
		ID           string
		Content      string
		Tags         []string
		RandomNumber int8 // yes, the types are different, but they will be converted
	}

	// third, parse it (field names must match exactly)
	var example Example

	err = parse.Parse(
		GRPCExample{
			ID:           "03973e64-358c-4a26-b095-150f18e8bfe7",
			Content:      "lorem ipsum",
			Tags:         []string{"foo", "bar"},
			RandomNumber: 9,
		},
		&example,

		// this function also accepts variadic options.
		// here we say that unknown fields will result parsing error (by default they won't. just like json)
		parse.WithDisallowUnknownFields(),
	)
	if err != nil {
		log.Fatalln(err)
	}

	// parsed values are already validated, therefore we can use it
	fmt.Println(example.ID.Get())
	// Output: 03973e64-358c-4a26-b095-150f18e8bfe7

	// By the way, parsing from maps and slices is also supported.
}
```

## Parsing

If needed, you can parse arbitrary types into your schemas through `parse` package.
See [parse example](./examples/parse/main.go) for more information.

Parsing gRPC messages is also supported, see [grpc parse example](./examples/parse-grpc/main.go)

## Performance

**TL;DR:** you can use codegen for max performance (0-1% overhead) or fallback to reflection (35% overhead).

This library does not affect unmarshalling performance itself.
You can expect it to be just as fast as a regular unmarshalling.

**However!**

Validation, by default, requires reflection to traverse over all struct fields. Again, reflection is only used to traverse fields, validators themself do not use reflection at all.

Such reflection traversal introduce ~35% performance overhead.

As an alternative, you can use [schemagen](./cmd/schemagen) [WIP] to generate field traversal logic.
As a result, overhead will reduced to 0-1% even for large structures. No need to change anything else.

Using reflection is easier because it does not require any codegen setup, but it does introduce minor performance decrease.

Unless performance is top-priority and validation is indeed a bottleneck (usually it's not), I'd recommend sticking with the reflection - it makes your codebase simpler to maintain.

**Benchmark:**

```
goos: darwin
goarch: arm64
pkg: github.com/metafates/schema/bench
cpu: Apple M3 Pro
BenchmarkUnmarshalJSON/reflection/with_validation-12      66221 ns/op
BenchmarkUnmarshalJSON/reflection/without_validation-12   45593 ns/op
BenchmarkUnmarshalJSON/codegen/with_validation-12         45936 ns/op
BenchmarkUnmarshalJSON/codegen/without_validation-12      45649 ns/op
```

## Validators

Not listed here yet, but can see a full list
of available validators in [validate/validators.go](./validate/validators.go)

## TODO

- [x] Support for manual construction (similar to `.parse(...)` in zod) (using codegen)
- [ ] Stabilize API
- [x] Better documentation
- [x] More tests
- [x] Improve performance. It should not be a bottleneck for most usecases, especially for basic CRUD apps. Still, there is a room for improvement!
- [x] Add benchmarks for validators itself. E.g. email validator
- [x] More validation types as seen in https://github.com/go-playground/validator
