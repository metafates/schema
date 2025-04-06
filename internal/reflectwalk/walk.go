package reflectwalk

import (
	"fmt"
	"reflect"
)

// FieldVisitor defines a function signature for the callback
// It receives the path to the field and its value
type FieldVisitor func(path string, value reflect.Value) error

// WalkFields traverses all fields in the given value, calling visitor for each field
func WalkFields(data any, visitor FieldVisitor) error {
	// Track visited pointers to avoid infinite recursion on cycles
	visited := make(map[uintptr]bool)
	return walkRecursive("", reflect.ValueOf(data), visitor, visited)
}

// walkRecursive recursively walks the reflected value
func walkRecursive(path string, v reflect.Value, visitor FieldVisitor, visited map[uintptr]bool) error {
	// Handle invalid values
	if !v.IsValid() {
		return visitor(path, v)
	}

	// First, call the visitor function on the current value
	if err := visitor(path, v); err != nil {
		return err
	}

	// Dereference pointers and unwrap interfaces
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		if v.IsNil() {
			return nil
		}

		// Check for cycles in pointers
		if v.Kind() == reflect.Ptr {
			ptr := v.Pointer()
			if visited[ptr] {
				return nil // Skip already visited pointers
			}
			visited[ptr] = true
		}

		return walkRecursive(path, v.Elem(), visitor, visited)
	}

	// Handle different value types
	switch v.Kind() {
	case reflect.Struct:
		// Walk struct fields
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			// Skip unexported fields that can't be interfaced with
			if !field.CanInterface() {
				continue
			}

			fieldName := t.Field(i).Name
			fieldPath := path
			if path != "" {
				fieldPath += "."
			}
			fieldPath += fieldName

			if err := walkRecursive(fieldPath, field, visitor, visited); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		// Skip nil slices after calling visitor
		if v.Kind() == reflect.Slice && v.IsNil() {
			return nil
		}

		// Walk array/slice elements
		for i := 0; i < v.Len(); i++ {
			indexPath := fmt.Sprintf("%s[%d]", path, i)
			if err := walkRecursive(indexPath, v.Index(i), visitor, visited); err != nil {
				return err
			}
		}

	case reflect.Map:
		// Skip nil maps after calling visitor
		if v.IsNil() {
			return nil
		}

		// Walk map keys and values
		for _, key := range v.MapKeys() {
			// Format key string representation for the path
			var keyStr string
			if key.Kind() == reflect.String {
				keyStr = key.String()
			} else {
				keyStr = fmt.Sprintf("%v", key.Interface())
			}

			// Walk map value with key path
			valuePath := fmt.Sprintf("%s[%s]", path, keyStr)
			if err := walkRecursive(valuePath, v.MapIndex(key), visitor, visited); err != nil {
				return err
			}
		}
	}

	return nil
}
