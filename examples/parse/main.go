package main

import (
	"fmt"
	"log"

	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/required"
)

func main() {
	type Dst struct {
		Name   required.NonZero[string]
		Age    int
		Arr    []struct{ X int }
		Nested struct{ Foo int }
	}

	for i, src := range []any{
		map[string]any{
			"Name":   "hi",
			"Age":    2,
			"Arr":    []map[string]any{{"X": 9}},
			"Nested": map[string]any{"Foo": 249},
		},
		struct {
			Name   string
			Age    int
			Arr    []struct{ X int }
			Nested struct{ Foo int }
		}{
			Name:   "hi",
			Age:    2,
			Arr:    []struct{ X int }{{X: 9}},
			Nested: struct{ Foo int }{Foo: 249},
		},
		struct{ Name string }{Name: "hello"},
	} {
		var dst Dst

		if err := parse.Parse(&dst, src); err != nil {
			log.Fatalf("%d: %v", i, err)
		}

		fmt.Printf("dst: %#+v\n", dst)
	}
}
