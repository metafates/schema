package wrap

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

var schemaType = reflect.TypeFor[interface{ IsSchema() }]()

// Wrap is a type that adds support for checking the presence of the required fields, if any
type Wrap[T any] struct{ Inner T }

func (w *Wrap[T]) UnmarshalJSON(data []byte) error {
	inner := reflect.TypeFor[T]()

	switch inner.Kind() {
	default:
		return json.Unmarshal(data, &w.Inner)

	case reflect.Slice, reflect.Array:
		panic("unimplemented")

	case reflect.Struct:
		// TODO: find a way to reduce calls to marshal/unmarshal

		var unmarshalled map[string]any

		if err := json.Unmarshal(data, &unmarshalled); err != nil {
			return err
		}

		fields, err := getRequiredFields(inner, "json")
		if err != nil {
			return err
		}

		mergeFields(fields, unmarshalled)

		filled, err := json.Marshal(unmarshalled)
		if err != nil {
			return err
		}

		return json.Unmarshal(filled, &w.Inner)
	}
}

func mergeFields(source, target map[string]any) {
	for k, expectedAtKey := range source {
		if _, ok := target[k]; ok {
			continue
		}

		if expectedAtKey == nil {
			target[k] = nil
			continue
		}

		target[k] = make(map[string]any)

		mergeFields(expectedAtKey.(map[string]any), target[k].(map[string]any))
	}
}

func getRequiredFields(t reflect.Type, forTag string) (map[string]any, error) {
	if t == nil || t.Implements(schemaType) {
		return nil, nil
	}

	switch t.Kind() {
	default:
		return nil, nil

	case reflect.Struct:
		fields := make(map[string]any)

		for i := 0; i < t.NumField(); i++ {
			field := t.Field(i)

			name := getName(field, forTag)

			if field.Type.Implements(schemaType) {
				if name == "-" {
					return nil, fmt.Errorf(`%s: "-" is used with required type`, field.Name)
				}

				fields[name] = nil
			} else if field.Type.Kind() == reflect.Struct {
				nested, err := getRequiredFields(field.Type, forTag)
				if err != nil {
					return nil, fmt.Errorf("%s: %w", field.Name, err)
				}
				fields[name] = nested
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
