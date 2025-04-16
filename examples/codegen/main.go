package main

import (
	"fmt"
	"log"
	"time"

	schemajson "github.com/metafates/schema/encoding/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate/charset"

	_ "embed"
)

//go:generate schemagen -type User

type User struct {
	ID    required.UUID[string]                          `json:"id"`
	Name  required.NonZeroCharset[string, charset.Print] `json:"name"`
	Birth optional.InPast[time.Time]                     `json:"birth"`

	Meta struct {
		Preferences optional.UniqueSlice[string] `json:"preferences"`
		Admin       bool                         `json:"admin"`
	} `json:"meta"`

	Friends []UserFriend `json:"friends"`

	Addresses []struct {
		Tag       optional.NonZeroCharset[string, charset.Print] `json:"tag"`
		Latitude  required.Latitude[float64]                     `json:"latitude"`
		Longitude required.Longitude[float64]                    `json:"longitude"`
	} `json:"addresses"`
}

type UserFriend struct {
	ID   required.UUID[string]                          `json:"id"`
	Name required.NonZeroCharset[string, charset.Print] `json:"name"`
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

		fmt.Println(schemajson.Unmarshal([]byte(invalidName), &user))
		// validate: .Name: string contains unprintable character
	}
}
