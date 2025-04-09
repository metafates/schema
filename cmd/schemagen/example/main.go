package main

import (
	"time"

	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
)

//go:generate schemagen -type=MyStruct

type MyStruct struct {
	Name  required.NonEmpty[string]  `json:"name"`
	Birth optional.InPast[time.Time] `json:"birth"`
	Anon  struct{ Foo required.ASCII[string] }
	Map   map[string]required.Any[string]
	Slice [][]map[string]required.Email[string]
	Bio   string
	Ptr   *Other
}

type Other struct {
	ID        required.NonEmpty[string]
	JustValue string
}
