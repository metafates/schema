// Code generated by schemagen; DO NOT EDIT.

package bench

import (
	"fmt"
	"github.com/metafates/schema/required"
	validate "github.com/metafates/schema/validate"
)

// Ensure that [DataWithGen] type was not changed
func _() {
	type locked []struct {
		ID         required.NonEmpty[string]   `json:"_id"`
		Index      int                         `json:"index"`
		GUID       string                      `json:"guid"`
		IsActive   bool                        `json:"isActive"`
		Balance    string                      `json:"balance"`
		Picture    string                      `json:"picture"`
		Age        required.Positive[int]      `json:"age"`
		EyeColor   string                      `json:"eyeColor"`
		Name       string                      `json:"name"`
		Gender     string                      `json:"gender"`
		Company    string                      `json:"company"`
		Email      string                      `json:"email"`
		Phone      string                      `json:"phone"`
		Address    string                      `json:"address"`
		About      string                      `json:"about"`
		Registered string                      `json:"registered"`
		Latitude   required.Latitude[float64]  `json:"latitude"`
		Longitude  required.Longitude[float64] `json:"longitude"`
		Tags       []string                    `json:"tags"`
		Friends    []struct {
			ID   int    `json:"id"`
			Name string `json:"name"`
		} `json:"friends"`
		Greeting      string `json:"greeting"`
		FavoriteFruit string `json:"favoriteFruit"`
	}
	var v DataWithGen
	// Compiler error signifies that the type definition have changed.
	// Re-run the schemagen command to regenerate this file.
	_ = locked(v)
}

// Validate implements [validate.Validateable]
func (x DataWithGen) Validate() error {
	for i0 := range x {
		{
			err0 := validate.Validate(&x[i0].ID)
			if err0 != nil {
				return validate.ValidationError{Inner: err0}.WithPath(fmt.Sprintf("[%v].ID", i0))
			}
			err1 := validate.Validate(&x[i0].Age)
			if err1 != nil {
				return validate.ValidationError{Inner: err1}.WithPath(fmt.Sprintf("[%v].Age", i0))
			}
			err2 := validate.Validate(&x[i0].Latitude)
			if err2 != nil {
				return validate.ValidationError{Inner: err2}.WithPath(fmt.Sprintf("[%v].Latitude", i0))
			}
			err3 := validate.Validate(&x[i0].Longitude)
			if err3 != nil {
				return validate.ValidationError{Inner: err3}.WithPath(fmt.Sprintf("[%v].Longitude", i0))
			}
		}
	}
	return nil
}
