package main

import (
	"fmt"
	"log"
	"time"

	schemajson "github.com/metafates/schema/encoding/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate/charset"
)

// Let's assume we have the following structure
type User struct {
	ID    required.UUID[string]
	Name  required.Charset[string, charset.Print]
	Birth optional.InPast[time.Time]

	FavoriteNumber int

	Friends []Friend
}

type Friend struct {
	ID   required.UUID[string]
	Name required.Charset[string, charset.Print]
}

func main() {
	// We can use json to unmarshal data to it
	{
		var user User

		const data = `
		{
			"ID": "2c376d16-321d-43b3-8648-2e64798cc6b3",
			"Name": "john",
			"FavoriteNumber": 42,
			"Friends": [
				{"ID": "7f735045-c8d2-4a60-9184-0fc033c40a6a", "Name": "jane"}
			]
		}
		`

		if err := schemajson.Unmarshal([]byte(data), &user); err != nil {
			log.Fatalln(err)
		}

		fmt.Println(user.Friends[0].ID.Get()) // 7f735045-c8d2-4a60-9184-0fc033c40a6a
	}

	// But we can also construct our user manually through parsing!
	//
	// For example, we can have some data as map.
	// Note, that field names should match go names exactly. Field tags (`json:"myField"`) won't be considered
	dataAsMap := map[string]any{
		"ID":             "2c376d16-321d-43b3-8648-2e64798cc6b3",
		"Name":           "john",
		"FavoriteNumber": 42,
		"Friends": []map[string]any{
			{"ID": "7f735045-c8d2-4a60-9184-0fc033c40a6a", "Name": "jane"},
		},
	}

	// structs are also supported. It is ok to omit non-required fields
	dataAsStruct := struct {
		ID             string
		Name           []byte // types will be converted, if possible. If not - error is returned
		FavoriteNumber uint16
	}{
		ID:             "2c376d16-321d-43b3-8648-2e64798cc6b3",
		Name:           []byte("john"),
		FavoriteNumber: 42,
	}

	// Now let's parse it
	for _, data := range []any{
		dataAsMap,
		dataAsStruct,
	} {
		var user User

		if err := parse.Parse(data, &user); err != nil {
			log.Fatalln(err)
		}

		fmt.Println(user.ID.Get()) // 2c376d16-321d-43b3-8648-2e64798cc6b3
	}
}
