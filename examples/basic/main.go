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

	data := []byte(`{"id": "hi", "email": "John Doe <john@example.com>"}`)

	if err := schemajson.Unmarshal(data, &r); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(r.ID.Value())
	fmt.Println(r.Email.Value())
}
