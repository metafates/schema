package main

import (
	"time"

	"github.com/metafates/schema/constraint"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
)

//go:generate schemagen -type User

type NonEmptyPrintable[T constraint.Text] struct {
	validate.And[
		T,
		validate.NonEmpty[T],
		validate.Printable[T],
	]
}

type User struct {
	ID    required.UUID[string]                              `json:"id"`
	Name  required.Custom[string, NonEmptyPrintable[string]] `json:"name"`
	Birth optional.InPast[time.Time]                         `json:"birth"`

	Meta struct {
		Preferences optional.Unique[[]string] `json:"preferences"`
		Admin       bool                      `json:"admin"`
	} `json:"meta"`

	// uuid -> name
	Friends map[required.UUID[string]]required.Custom[string, NonEmptyPrintable[string]] `json:"friends"`
}
