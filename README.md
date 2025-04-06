# Schema

> Work in progress!

Type-safe schema guarded structs for Go with generics and a bit of magic.

No stable version yet, but you can use it like that.

```bash
go get github.com/metafates/schema@main
```

**Work in progress, API may change significantly without further notice!**

## Example

See [examples](./examples) for more examples

```go
package main

import (
	"fmt"
	"log"

	schemajson "github.com/metafates/schema/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
)

type Request struct {
    // this field is required and will result unmarshal error if missing
	ID required.NonEmpty[string] `json:"id"`

    // this field is optional and will be empty if missing.
    // if it is present in json it will be validated for email syntax.
    // if validation fails, unmarshal will return error
	Email optional.Email[string] `json:"email"`

    // this field is just a regular go value, which will be an empty string if missing.
    // no validation or further logic is attached to it
	Comment string
}

func main() {
	var r Request

	data := []byte(`{"id":"hi", "email":"John Doe <john@example.com>"}`)

	if err := schemajson.Unmarshal(data, &r); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(r.ID.Value())
	fmt.Println(r.Email.Value())
}
```

## Custom types

It is possible to create custom optional or required type.
See [examples/custom](./examples/custom/main.go) for an advanced example

```go
// You can create your custom validator like that
type ShortStr struct{}

func (ShortStr) Validate(v string) error {
	if len(v) > 10 {
		return errors.New("string is too long")
	}

	return nil
}

// Or combine existing validators using special [validate.Combined] type
type NonZeroPositive[T constraint.Real] struct {
	validate.Combined[
		validate.NonEmpty[T],
		validate.Positive[T],
		T,
	]
}

// Now you can use you validators like that
type Request struct {
	Foo required.Custom[string, ShortStr]                  `json:"foo"`
	Bar optional.Custom[float64, NonZeroPositive[float64]] `json:"bar"`
}

// You can also create an alias for such custom types like that
// In fact, this is how builtin types for required and optional are implemented
type RequiredShortString struct {
	required.Custom[string, ShortStr]
}

type Request2 struct {
	Foo RequiredShortString `json:"foo"`
}
```
