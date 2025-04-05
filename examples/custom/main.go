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

type ShortStr struct{}

func (ShortStr) Validate(v string) error {
	if len(v) > 10 {
		return errors.New("string is too long")
	}

	return nil
}

type NonZeroPositive[T constraint.Real] struct {
	validate.Combined[
		validate.NonEmpty[T],
		validate.Positive[T],
		T,
	]
}

type Request struct {
	Foo required.Custom[string, ShortStr]                  `json:"foo"`
	Bar optional.Custom[float64, NonZeroPositive[float64]] `json:"bar"`
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
