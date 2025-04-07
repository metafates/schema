package main

import (
	"encoding/json"
	"fmt"
	"log"

	schemajson "github.com/metafates/schema/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/wrap"
)

type Address struct {
	Lat required.Latitude[float64]  `json:"lat"`
	Lon required.Longitude[float64] `json:"lon"`
}

type Request struct {
	// this field is required and will result unmarshal error if missing
	ID required.NonEmpty[string] `json:"id"`

	// this field is optional and will be empty if missing.
	// if it is present in json it will be validated for email syntax.
	// if validation fails, unmarshal will return error
	Email optional.Email[string] `json:"email"`

	Address required.Any[Address] `json:"address"`

	// this field is just a regular go value, which will be an empty string if missing.
	// no validation or further logic is attached to it
	Comment string
}

func main() {
	var r Request

	data := []byte(`{"id": "hi", "email": "John Doe <john@example.com>", "address": {"lat":24,"lon": 91}}`)

	// this is important! validation won't work otherwise.
	// you need to use [schemajson.Unmarshal]
	{
		if err := schemajson.Unmarshal(data, &r); err != nil {
			log.Fatalln(err)
		}
	}

	// as an alternative you can also do something like that.
	// in fact, this is exactly how [schemajson.Unmarshal] is implemented.
	{
		var wrapped wrap.Wrapped[Request]

		if err := json.Unmarshal(data, &wrapped); err != nil {
			log.Fatalln(err)
		}

		r = wrapped.Inner
	}

	fmt.Println(r.ID.Value())
	fmt.Println(r.Email.Value())
}
