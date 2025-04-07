# 📐 Schema

> Work in progress!

Type-safe schema guarded structs for Go with generics and a bit of magic.

No stable version yet, but you can use it like that.

```bash
go get github.com/metafates/schema@main
```

**Work in progress, API may change significantly without further notice! This is just an experiment for now**

## Example

See [examples](./examples) for more examples

```go
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	schemajson "github.com/metafates/schema/json"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
	"github.com/metafates/schema/validate"
	"github.com/metafates/schema/validate/wrap"
)

// Let's assume we have a request which accepts an user
type User struct {
	// User name is required and must not be empty
	Name required.NonEmpty[string] `json:"name"`

	// Birth date is optional, which means it could be null.
	// However, if passed, it must be an any valid [time.Time]
	Birth optional.Any[time.Time] `json:"birth"`

	// Same for email. It is optional, therefore it could be null.
	// But, if passed, not only it must be a valid string, but also a valid RFC 5322 email string
	Email optional.Email[string] `json:"email"`

	// Bio is just a regular string. It may be empty, may be not.
	// No further logic is attached to it.
	Bio string `json:"bio"`
}

// We could also have an address
type Address struct {
	// Latitude is required and must be a valid latitude (range [-90; 90])
	Latitude required.Latitude[float64] `json:"latitude"`

	// Longitude is also required and must be a valid longitude (range [-180; 180])
	Longitude required.Longitude[float64] `json:"longitude"`
}

// But wait, what about custom types?
// We might want (for some reason) a field which accepts only short strings (<10 bytes).
// Let's see how we might implement it.

// let's define a custom validator for short strings.
// it should not contain any fields in itself, they won't be initialized or used in any way.
type ShortStr struct{}

// this function implements a special validator interface
func (ShortStr) Validate(v string) error {
	if len(v) >= 10 {
		return errors.New("string is too long")
	}

	return nil
}

// that's it, basically. now we can use this validator in our request.

// but we can go extreme! we can combine multiple validators using types
type ASCIIShortStr struct {
	// both ASCII and ShortStr must be satisfied.
	// you can also use [validate.Or] to ensure that at least one condition is satisfied.
	validate.And[
		validate.ASCII[string],
		ShortStr,
		string,
	]
}

// Now, our final request may look something like that
type Request struct {
	// User is required by default.
	// Because, if not passed, it will be empty.
	// Therefore required fields in user will also be empty which will result a missing fields error.
	User User `json:"user"`

	// Address is, however, optional. It could be null.
	// But, if passed, it must be a valid address with respect to its fields validation (required lat/lon)
	Address optional.Any[Address] `json:"address"`

	// this is how we can use our validator in custom type.
	// we could make an alias for that custom required type, if needed.
	MyShortString required.Custom[string, ShortStr] `json:"myShortString"`

	// same as for [Request.MyShortString] but using an optional instead.
	ASCIIShortString optional.Custom[string, ASCIIShortStr] `json:"asciiShortString"`
}

func main() {
	// here we create a VALID json for our request. we will try to pass an invalid one later.
	data := []byte(`
{
	"user": {
		"name": "john",
		"email": "john@example.com (comment)",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "foo"
}`)

	var request Request

	// in order to make it all work, we need to use a special unmarshal function [schemajson.Unmarshal].
	if err := schemajson.Unmarshal(data, &request); err != nil {
		log.Fatalln(err)
	}

	// you can also use regular [json.Unmarshal] if you need to.
	{
		// just wrap your type in a special wrapper.
		// only do it once for root type, no need to wrap each field
		var wrapped wrap.Wrapped[Request]

		// will also work for [json.Decoder]
		if err := json.Unmarshal(data, &wrapped); err != nil {
			log.Fatalln(err)
		}

		// and unwrap it
		request = wrapped.Inner
	}

	// moreover, it possible to separate validation from unmarshalling
	{
		// unmarshal it any way you want, not only json
		if err := json.Unmarshal(data, &request); err != nil {
			log.Fatalln(err)
		}

		// and validate it manually, that's it
		if err := validate.Recursively(request); err != nil {
			log.Fatalln(err)
		}

		// in fact, this is what [wrap.Wrapped] does.
		// this seems to be more flexible approach, the API change to emphasize it
	}

	// now that we have successfully unmarshalled our json, we can use request fields.
	// to access values of our schema-guarded fields we can use .Value() method
	//
	// NOTE: calling this method for required types BEFORE we have
	// unmarshalled our request will panic intentionally.
	fmt.Println(request.User.Name.Value()) // output: john

	// optional values return a tuple: a value and boolean stating its presence
	email, ok := request.User.Email.Value()
	fmt.Println(email, ok) // output: john@example.com true

	// birth is missing
	birth, ok := request.User.Birth.Value()
	fmt.Println(birth, ok) // output: 0001-01-01 00:00:00 +0000 UTC false

	// let's try to pass an INVALID jsons.
	invalidEmail := []byte(`
{
	"user": {
		"name": "john",
		"email": "john@@@example.com",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "foo"
}`)

	invalidShortStr := []byte(`
{
	"user": {
		"name": "john",
		"email": "john@example.com",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "super long string that shall not pass!!!!!!!!"
}`)

	missingUserName := []byte(`
{
	"user": {
		"email": "john@example.com",
		"bio": "lorem ipsum"
	},
	"address": {
		"latitude": 81.111,
		"longitude": 100.101
	},
	"myShortString": "foo"
}`)

	fmt.Println(schemajson.Unmarshal(invalidEmail, &request))
	// validate: User.Email: mail: missing '@' or angle-addr

	fmt.Println(schemajson.Unmarshal(invalidShortStr, &request))
	// validate: MyShortString: string is too long

	fmt.Println(schemajson.Unmarshal(missingUserName, &request))
	// validate: User.Name: missing value

	// You can also check if it was validation error or any other json error
	err := schemajson.Unmarshal(missingUserName, &request)

	var validationErr schemaerror.ValidationError
	if errors.As(err, &validationErr) {
		fmt.Println("validation error occured:", validationErr)
	}
}
```

## Performance

This library does not affect unmarshalling performance itself. You can expect it to be just as fast as a regular unmarshalling. **However**, the validation itself does have an overhead:

[Benchmark source code](./validate/wrap/wrap_test.go)

```
goos: darwin
goarch: arm64
pkg: github.com/metafates/schema/validate/wrap
cpu: Apple M3 Pro
BenchmarkUnmarshalJSON/unmarshal_and_validation_with_wrap-12               12141             98618 ns/op
BenchmarkUnmarshalJSON/separate_unmarshal_and_validation_manually-12       18802             63793 ns/op
BenchmarkUnmarshalJSON/unmarshal_without_validation-12                     27460             43477 ns/op
```

As you can see, it's ~40-100% performance slowdown when using validation. The reason for it is that validation requires iterating over all unmarshalled fields and validate required and optional fields. There's a room for improvement!

## TODO

- [ ] Better documentation
- [ ] More tests
- [ ] Improve performance. It should not be a bottleneck for most usecases, especially for basic CRUD apps. Still, there is a room for improvement!
- [ ] More validation types as seen in https://github.com/go-playground/validator
