package main

import (
	"fmt"
	"log"
	"time"

	schemajson "github.com/metafates/schema/encoding/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"

	_ "embed"
)

//go:generate schemagen -type User

type User struct {
	ID    required.UUID[string]              `json:"id"`
	Name  required.NonEmptyPrintable[string] `json:"name"`
	Birth optional.InPast[time.Time]         `json:"birth"`

	Meta struct {
		Preferences optional.Unique[[]string, string] `json:"preferences"`
		Admin       bool                              `json:"admin"`
	} `json:"meta"`

	Friends []UserFriend `json:"friends"`

	Addresses []struct {
		Tag       optional.NonEmptyPrintable[string] `json:"tag"`
		Latitude  required.Latitude[float64]         `json:"latitude"`
		Longitude required.Longitude[float64]        `json:"longitude"`
	} `json:"addresses"`
}

type UserFriend struct {
	ID   required.UUID[string]              `json:"id"`
	Name required.NonEmptyPrintable[string] `json:"name"`
}

var (
	//go:embed valid_data.json
	validData string

	//go:embed invalid_data.json
	invalidName string
)

func main() {
	{
		var user User

		if err := schemajson.Unmarshal([]byte(validData), &user); err != nil {
			log.Fatalln(err)
		}

		fmt.Println(user.ID.Get())
		// 32541419-1294-47e4-b070-833db7684866
	}
	{
		var user User

		if err := schemajson.Unmarshal([]byte(invalidName), &user); err != nil {
			log.Fatalln(err)
			// validate: .Name: string contains unprintable character
		}

		fmt.Println(user.ID.Get())
	}
}
