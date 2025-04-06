package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/metafates/schema/constraint"
	schemajson "github.com/metafates/schema/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
)

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

func main() {
	data := []byte(`
{
	"foo": "short str",
	"bar": 42.2
}
`)

	var r Request

	if err := schemajson.Unmarshal(data, &r); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(r.Foo.Value())
	fmt.Println(r.Bar.Value())
}
