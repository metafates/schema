package main

import (
	"fmt"
	"log"
	"time"

	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate/charset"
)

type User struct {
	Name required.Charset[string, charset.Print]

	Birth optional.InPast[time.Time]

	Bio string
}

func main() {
	{
		var user User

		err := parse.Parse(&user, map[string]any{
			"Name":  "john",
			"Birth": time.Date(1900, time.September, 10, 0, 0, 0, 0, time.UTC),
			"Bio":   "...",
		})
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(user.Name.Get())
		// john
	}

	{
		var user User

		err := parse.Parse(&user, struct{ Name string }{Name: "john"})
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Println(user.Name.Get())
		// john
	}
}
