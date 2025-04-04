package main

import (
	"fmt"
	"log"

	"github.com/metafates/schema"
	rjson "github.com/metafates/schema/json"
)

type Request struct {
	ID schema.RequiredNonEmpty[string] `json:"id"`

	Bet schema.Required[float64] `json:"bet"`

	Comment string
}

func main() {
	var r Request

	if err := rjson.Unmarshal([]byte(`{"bet":42.42, "id":"hi"}`), &r); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(r.ID.Value())
}
