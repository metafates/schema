package reflectwalk

import (
	"fmt"
	"reflect"
	"strconv"
)

// FieldVisitor defines a function signature for the callback.
// It receives the path to the field and its value.
type FieldVisitor func(path string, value reflect.Value) error

// WalkFields traverses all fields in the given value, calling visitor for each field.
func WalkFields(data any, visitor FieldVisitor) error {
	// Track visited pointers to prevent infinite recursion on cycles.
	visited := make(map[uintptr]bool)
	return walkRecursive("", reflect.ValueOf(data), visitor, visited)
}

func walkRecursive(path string, v reflect.Value, visitor FieldVisitor, visited map[uintptr]bool) error {
	// If we have an invalid (zero) reflect.Value, just fire the visitor.
	if !v.IsValid() {
		return visitor(path, v)
	}

	// Call the visitor on the current value first.
	if err := visitor(path, v); err != nil {
		return err
	}

	kind := v.Kind()

	switch kind {
	case reflect.Ptr:
		// Check for nil pointer and visited pointer cycle first.
		if v.IsNil() {
			return nil
		}
		ptr := v.Pointer()
		if visited[ptr] {
			return nil
		}
		visited[ptr] = true
		return walkRecursive(path, v.Elem(), visitor, visited)

	case reflect.Interface:
		// If interface is nil, nothing to do.
		if v.IsNil() {
			return nil
		}
		return walkRecursive(path, v.Elem(), visitor, visited)

	case reflect.Struct:
		// Walk struct fields.
		t := v.Type()
		numField := v.NumField()
		for i := 0; i < numField; i++ {
			fieldVal := v.Field(i)
			// Skip unexported fields.
			if !fieldVal.CanInterface() {
				continue
			}

			fieldName := t.Field(i).Name

			var fieldPath string
			if path == "" {
				fieldPath = "." + fieldName
			} else {
				fieldPath = path + "." + fieldName
			}

			if err := walkRecursive(fieldPath, fieldVal, visitor, visited); err != nil {
				return err
			}
		}

	case reflect.Array, reflect.Slice:
		// After visiting the slice itself, skip if it's nil.
		if kind == reflect.Slice && v.IsNil() {
			return nil
		}

		length := v.Len()
		for i := 0; i < length; i++ {
			indexPath := path + "[" + strconv.Itoa(i) + "]"
			if err := walkRecursive(indexPath, v.Index(i), visitor, visited); err != nil {
				return err
			}
		}

	case reflect.Map:
		// After visiting the map itself, skip if it's nil.
		if v.IsNil() {
			return nil
		}

		keys := v.MapKeys()
		for _, key := range keys {
			// Convert map key to string inlined for performance.
			var keyStr string

			switch k := key.Kind(); k {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				keyStr = strconv.FormatInt(key.Int(), 10)

			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
				keyStr = strconv.FormatUint(key.Uint(), 10)

			case reflect.Bool:
				if key.Bool() {
					keyStr = "true"
				} else {
					keyStr = "false"
				}

			case reflect.Float32, reflect.Float64:
				keyStr = strconv.FormatFloat(key.Float(), 'g', -1, 64)

			case reflect.String:
				keyStr = key.String()

			default:
				// Fallback for complex or otherwise unsupported key kinds.
				keyStr = fmt.Sprint(key.Interface())
			}

			valuePath := path + "[" + keyStr + "]"
			val := v.MapIndex(key)
			if err := walkRecursive(valuePath, val, visitor, visited); err != nil {
				return err
			}
		}
	}

	return nil
}
