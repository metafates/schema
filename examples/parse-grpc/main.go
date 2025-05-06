package main

import (
	"fmt"
	"log"

	"github.com/metafates/schema/examples/parse-grpc/pb"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/parse"
	"github.com/metafates/schema/required"
)

type AddressBook struct {
	People []Person
}

type Person struct {
	Name   required.NonZero[string]
	Id     required.Positive0[int32]
	Email  optional.Email[string]
	Phones []PhoneNumber
}

type PhoneType int

const (
	PhoneTypeMobile = iota
	PhoneTypeHome
	PhoneTypeWork
)

type PhoneNumber struct {
	Type   optional.Any[PhoneType]
	Number required.NonZero[string]
}

func main() {
	options := []parse.Option{
		parse.WithDisallowUnknownFields(),
	}

	// let's parse valid address book from grpc
	{
		var book AddressBook

		err := parse.Parse(pb.AddressBook{
			People: []*pb.Person{
				{
					Name:  "Example Name",
					Id:    12345,
					Email: "name@example.com",
					Phones: []*pb.Person_PhoneNumber{
						{
							Number: "123-456-7890",
							Type:   pb.Person_HOME,
						},
						{
							Number: "222-222-2222",
							Type:   pb.Person_MOBILE,
						},
						{
							Number: "111-111-1111",
							Type:   pb.Person_WORK,
						},
					},
				},
			},
		}, &book, options...)
		if err != nil {
			log.Fatalln(err)
		}

		fmt.Printf("book.People: %v\n", len(book.People))
		// 1

		fmt.Printf("book.People[0].Name: %v\n", book.People[0].Name.Get())
		// Example Name

		fmt.Printf("book.People[0].Phones: %v\n", len(book.People[0].Phones))
		// 3

		fmt.Printf(
			"pb.Person_MOBILE == PhoneTypeMobile = %v\n",
			book.People[0].Phones[1].Type.Must() == PhoneTypeMobile,
		)
		// pb.Person_MOBILE == PhoneTypeMobile = true
	}

	// now let's try to trigger error by violating the schema
	{
		var book AddressBook

		err := parse.Parse(pb.AddressBook{
			People: []*pb.Person{
				{
					Name:  "Example Name",
					Id:    12345,
					Email: "not a valid email",
					Phones: []*pb.Person_PhoneNumber{
						{
							Number: "123-456-7890",
							Type:   pb.Person_HOME,
						},
						{
							Type: pb.Person_MOBILE,
						},
						{
							Number: "111-111-1111",
							Type:   pb.Person_WORK,
						},
					},
				},
			},
		}, &book, options...)

		fmt.Println(err) // [0].Email: mail: no angle-addr
	}
}
