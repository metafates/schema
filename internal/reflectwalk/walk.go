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

	switch v.Kind() {
	case reflect.Ptr:
		return walkPtr(path, v, visitor, visited)

	case reflect.Interface:
		return walkInterface(path, v, visitor, visited)

	case reflect.Struct:
		return walkStruct(path, v, visitor, visited)

	case reflect.Array, reflect.Slice:
		return walkSlice(path, v, visitor, visited)

	case reflect.Map:
		return walkMap(path, v, visitor, visited)

	default:
		return nil
	}
}

func walkPtr(path string, v reflect.Value, visitor FieldVisitor, visited map[uintptr]bool) error {
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
}

func walkInterface(path string, v reflect.Value, visitor FieldVisitor, visited map[uintptr]bool) error {
	// If interface is nil, nothing to do.
	if v.IsNil() {
		return nil
	}

	return walkRecursive(path, v.Elem(), visitor, visited)
}

func walkStruct(path string, v reflect.Value, visitor FieldVisitor, visited map[uintptr]bool) error {
	t := v.Type()

	for i := range v.NumField() {
		fieldVal := v.Field(i)

		// Skip unexported fields.
		if !fieldVal.CanInterface() {
			continue
		}

		fieldName := t.Field(i).Name
		fieldPath := path + "." + fieldName

		if err := walkRecursive(fieldPath, fieldVal, visitor, visited); err != nil {
			return err
		}
	}

	return nil
}

func walkSlice(path string, v reflect.Value, visitor FieldVisitor, visited map[uintptr]bool) error {
	// After visiting the slice itself, skip if it's nil.
	if v.Kind() == reflect.Slice && v.IsNil() {
		return nil
	}

	for i := range v.Len() {
		indexPath := path + "[" + strconv.Itoa(i) + "]"
		if err := walkRecursive(indexPath, v.Index(i), visitor, visited); err != nil {
			return err
		}
	}

	return nil
}

func walkMap(path string, v reflect.Value, visitor FieldVisitor, visited map[uintptr]bool) error {
	// After visiting the map itself, skip if it's nil.
	if v.IsNil() {
		return nil
	}

	keys := v.MapKeys()
	for _, key := range keys {
		valuePath := path + "[" + formatStr(key) + "]"
		val := v.MapIndex(key)

		if err := walkRecursive(valuePath, val, visitor, visited); err != nil {
			return err
		}
	}

	return nil
}

func formatStr(v reflect.Value) string {
	switch v.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(v.Int(), 10)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return strconv.FormatUint(v.Uint(), 10)

	case reflect.Bool:
		if v.Bool() {
			return "true"
		}

		return "false"

	case reflect.Float32, reflect.Float64:
		return strconv.FormatFloat(v.Float(), 'g', -1, 64)

	case reflect.String:
		return v.String()

	default:
		// Fallback for complex or otherwise unsupported key kinds.
		return fmt.Sprint(v.Interface())
	}
}
