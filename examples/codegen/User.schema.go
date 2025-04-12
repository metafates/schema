// Code generated by schemagen; DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
	validate "github.com/metafates/schema/validate"
	"time"
)

// Ensure that [User] type was not changed
func _() {
	type locked struct {
		ID    required.UUID[string]              `json:"id"`
		Name  required.NonEmptyPrintable[string] `json:"name"`
		Birth optional.InPast[time.Time]         `json:"birth"`
		Meta  struct {
			Preferences optional.Unique[[]string, string] `json:"preferences"`
			Admin       bool                              `json:"admin"`
		} `json:"meta"`
		Friends   []UserFriend `json:"friends"`
		Addresses []struct {
			Tag       optional.NonEmptyPrintable[string] `json:"tag"`
			Latitude  required.Latitude[float64]         `json:"latitude"`
			Longitude required.Longitude[float64]        `json:"longitude"`
		} `json:"addresses"`
	}
	var v User
	// Compiler error signifies that the type definition have changed.
	// Re-run the schemagen command to regenerate this file.
	_ = locked(v)
}

// Validate implements [validate.Validateable]
func (x *User) Validate() error {
	err0 := validate.Validate(&x.ID)
	if err0 != nil {
		return validate.ValidationError{Inner: err0}.WithPath(fmt.Sprintf(".ID"))
	}
	err1 := validate.Validate(&x.Name)
	if err1 != nil {
		return validate.ValidationError{Inner: err1}.WithPath(fmt.Sprintf(".Name"))
	}
	err2 := validate.Validate(&x.Birth)
	if err2 != nil {
		return validate.ValidationError{Inner: err2}.WithPath(fmt.Sprintf(".Birth"))
	}
	err3 := validate.Validate(&x.Meta.Preferences)
	if err3 != nil {
		return validate.ValidationError{Inner: err3}.WithPath(fmt.Sprintf(".Meta.Preferences"))
	}
	for i0 := range x.Friends {
		{
			err4 := validate.Validate(&x.Friends[i0])
			if err4 != nil {
				return validate.ValidationError{Inner: err4}.WithPath(fmt.Sprintf(".Friends[%v]", i0))
			}
		}
	}
	for i1 := range x.Addresses {
		{
			err5 := validate.Validate(&x.Addresses[i1].Tag)
			if err5 != nil {
				return validate.ValidationError{Inner: err5}.WithPath(fmt.Sprintf(".Addresses[%v].Tag", i1))
			}
			err6 := validate.Validate(&x.Addresses[i1].Latitude)
			if err6 != nil {
				return validate.ValidationError{Inner: err6}.WithPath(fmt.Sprintf(".Addresses[%v].Latitude", i1))
			}
			err7 := validate.Validate(&x.Addresses[i1].Longitude)
			if err7 != nil {
				return validate.ValidationError{Inner: err7}.WithPath(fmt.Sprintf(".Addresses[%v].Longitude", i1))
			}
		}
	}
	return nil
}
