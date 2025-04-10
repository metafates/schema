// Code generated by "schemagen -type MyStruct,ASlice,AMap,Basic,ABasicNested". DO NOT EDIT.

package main

import (
	"fmt"
	"github.com/metafates/schema/optional"
	"github.com/metafates/schema/required"
	validate "github.com/metafates/schema/validate"
	"time"
)

// Ensure types are not changed
func _() {
	type locked struct {
		Name  required.NonEmpty[string]  `json:"name"`
		Birth optional.InPast[time.Time] `json:"birth"`
		Anon  struct {
			Foo required.ASCII[string]
		}
		Map   map[string]required.Any[string]
		Slice [][]map[string]required.Email[string]
		Bio   string
		Ptr   *Other
		Ptr2  *[]string
	}
	// Compiler error signifies that the type definition have changed.
	// Re-run the schemagen command to regenerate validators.
	_ = locked(MyStruct{})
}

// Validate implementes [validate.Validateable]
func (x *MyStruct) Validate() error {
	v2 := &x.Name
	err1 := validate.Validate(v2)
	if err1 != nil {
		return validate.ValidationError{Inner: err1}.WithPath(fmt.Sprintf(".Name"))
	}
	v4 := &x.Birth
	err3 := validate.Validate(v4)
	if err3 != nil {
		return validate.ValidationError{Inner: err3}.WithPath(fmt.Sprintf(".Birth"))
	}
	v6 := &x.Anon.Foo
	err5 := validate.Validate(v6)
	if err5 != nil {
		return validate.ValidationError{Inner: err5}.WithPath(fmt.Sprintf(".Anon.Foo"))
	}
	for k7 := range x.Map {
		{
			v9 := x.Map[k7]
			err8 := validate.Validate(v9)
			x.Map[k7] = v9
			if err8 != nil {
				return validate.ValidationError{Inner: err8}.WithPath(fmt.Sprintf(".Map[%v]", k7))
			}
		}
	}
	for i10 := range x.Slice {
		{
			for i11 := range x.Slice[i10] {
				{
					for k12 := range x.Slice[i10][i11] {
						{
							v14 := x.Slice[i10][i11][k12]
							err13 := validate.Validate(v14)
							x.Slice[i10][i11][k12] = v14
							if err13 != nil {
								return validate.ValidationError{Inner: err13}.WithPath(fmt.Sprintf(".Slice[%v][%v][%v]", i10, i11, k12))
							}
						}
					}
				}
			}
		}
	}
	if x.Ptr != nil {
		{
			v16 := x.Ptr
			err15 := validate.Validate(v16)
			if err15 != nil {
				return validate.ValidationError{Inner: err15}.WithPath(fmt.Sprintf(".Ptr"))
			}
		}
	}
	return nil
}

// Ensure types are not changed
func _() {
	type locked []MyStruct
	// Compiler error signifies that the type definition have changed.
	// Re-run the schemagen command to regenerate validators.
	_ = locked(ASlice{})
}

// Validate implementes [validate.Validateable]
func (x ASlice) Validate() error {
	for i18 := range x {
		{
			v20 := &x[i18]
			err19 := validate.Validate(v20)
			if err19 != nil {
				return validate.ValidationError{Inner: err19}.WithPath(fmt.Sprintf("[%v]", i18))
			}
		}
	}
	return nil
}

// Ensure types are not changed
func _() {
	type locked map[string]MyStruct
	// Compiler error signifies that the type definition have changed.
	// Re-run the schemagen command to regenerate validators.
	_ = locked(AMap{})
}

// Validate implementes [validate.Validateable]
func (x AMap) Validate() error {
	for k21 := range x {
		{
			v23 := x[k21]
			err22 := validate.Validate(v23)
			x[k21] = v23
			if err22 != nil {
				return validate.ValidationError{Inner: err22}.WithPath(fmt.Sprintf("[%v]", k21))
			}
		}
	}
	return nil
}

// Ensure types are not changed
func _() {}

// Validate implementes [validate.Validateable]
func (x *Basic) Validate() error {
	return nil
}

// Ensure types are not changed
func _() {
	type locked map[string]map[string][]string
	// Compiler error signifies that the type definition have changed.
	// Re-run the schemagen command to regenerate validators.
	_ = locked(ABasicNested{})
}

// Validate implementes [validate.Validateable]
func (x ABasicNested) Validate() error {
	return nil
}
