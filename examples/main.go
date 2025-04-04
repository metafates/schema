package main

import (
	"fmt"
	"log"

	"github.com/metafates/required"
	requiredjson "github.com/metafates/required/json"
)

type Request struct {
	ID required.NonEmpty[string] `json:"id"`

	Bet required.T[float64] `json:"bet"`

	X struct {
		Nested required.Positive[int] `json:"nested"`
	} `json:"x"`

	Comment string
}

func main() {
	var r Request

	if err := requiredjson.Unmarshal([]byte(`{"bet":42.42, "id":"hi", "x": {"nested": 24}}`), &r); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(r.ID.Value())
}
