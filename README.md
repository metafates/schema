# Schema

> Work in progress!

Type-safe schema guarded structs for Go with generics and a bit of magic.

## Example

See [examples](./examples) for more examples

```go
package main

import (
	"fmt"
	"log"

	schemajson "github.com/metafates/schema/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
)

type Request struct {
	ID    required.NonEmpty[string] `json:"id"`
	Email optional.Email[string]    `json:"email"`

	Comment string
}

func main() {
	var r Request

	data := []byte(`{"id":"hi", "email":"John Doe <john@example.com>"}`)

	if err := schemajson.Unmarshal(data, &r); err != nil {
		log.Fatalln(err)
	}

	fmt.Println(r.ID.Value())
	fmt.Println(r.Email.Value())
}
```
