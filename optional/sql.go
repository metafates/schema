package optional

import (
	"database/sql"
	"database/sql/driver"
	"fmt"

	"github.com/metafates/schema/validate"
)

var _ interface {
	sql.Scanner
	driver.Valuer
} = (*Custom[any, validate.Validator[any]])(nil)

// Scan implements the [sql.Scanner] interface.
func (c *Custom[T, V]) Scan(src any) error {
	if src == nil {
		*c = Custom[T, V]{}
		return nil
	}

	var value T

	if scanner, ok := any(&value).(sql.Scanner); ok {
		if err := scanner.Scan(src); err != nil {
			return err
		}

		*c = Custom[T, V]{value: value, hasValue: true}
		return nil
	}

	if converted, err := driver.DefaultParameterConverter.ConvertValue(src); err == nil {
		if v, ok := converted.(T); ok {
			*c = Custom[T, V]{value: v, hasValue: true}
			return nil
		}
	}

	var nullable sql.Null[T]

	if err := nullable.Scan(src); err != nil {
		return err
	}

	if nullable.Valid {
		*c = Custom[T, V]{value: nullable.V, hasValue: true}
	} else {
		*c = Custom[T, V]{}
	}

	return nil
}

// Value implements the [driver.Valuer] interface.
//
// Use [Custom.Get] method instead for getting the go value
func (c Custom[T, V]) Value() (driver.Value, error) {
	if c.hasValue && !c.validated {
		panic("called UnmarshalJSON() on unvalidated value")
	}

	if !c.hasValue {
		return nil, nil
	}

	value, err := driver.DefaultParameterConverter.ConvertValue(c.value)
	if err != nil {
		return nil, fmt.Errorf("convert: %w", err)
	}

	return value, nil
}
