package wrap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var requiredType = reflect.TypeFor[interface{ IsRequired() }]()

type Wrap[T any] struct{ Inner T }

func (w *Wrap[T]) UnmarshalJSON(data []byte) error {
	inner := reflect.TypeFor[T]()

	switch inner.Kind() {
	default:
		return json.Unmarshal(data, &w.Inner)

	case reflect.Struct:
		var unmarshalled map[string]any

		if err := json.Unmarshal(data, &unmarshalled); err != nil {
			return err
		}

		fields, err := requiredFields(inner, "json")
		if err != nil {
			return err
		}

		propogateFields(fields, unmarshalled)

		filled, err := json.Marshal(unmarshalled)
		if err != nil {
			return err
		}

		return json.Unmarshal(filled, &w.Inner)
	}
}

func propogateFields(expected, actual map[string]any) {
	for k, expectedAtKey := range expected {
		if _, ok := actual[k]; ok {
			continue
		}

		if expectedAtKey == nil {
			actual[k] = nil
			continue
		}

		actual[k] = make(map[string]any)

		propogateFields(expectedAtKey.(map[string]any), actual[k].(map[string]any))
	}
}

func requiredFields(t reflect.Type, forTag string) (map[string]any, error) {
	if t == nil || t.Implements(requiredType) {
		return nil, nil
	}

	switch t.Kind() {
	default:
		return nil, nil

	case reflect.Struct:
		fields := make(map[string]any)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			nested, err := requiredFields(field.Type, forTag)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", field.Name, err)
			}

			// hack: we could check for sealed interface implementation
			// but it won't work across packages
			if field.Type.PkgPath() != "github.com/metafates/schema" {
				continue
			}

			name := getName(field, forTag)

			if name == "-" {
				return nil, fmt.Errorf(`%s: "-" is used with required type`, field.Name)
			}

			if nested != nil {
				fields[name] = nested
			} else {
				fields[name] = nil
			}
		}

		return fields, nil
	}
}

func getName(field reflect.StructField, tag string) string {
	tagValue := field.Tag.Get(tag)

	if tagValue == "" {
		return field.Name
	}

	values := strings.Split(field.Tag.Get(tag), ",")

	if len(values) == 0 {
		return field.Name
	}

	return values[0]
}
