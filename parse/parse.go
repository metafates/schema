package parse

import (
	"fmt"
	"reflect"
	"strconv"
)

type Parser interface {
	Parse(v any) error
}

// Parse fills the struct pointed to by src based on fields/keys in target.
// It returns an error if types cannot be assigned/converted.
func Parse(dst, src any) error {
	if parser, ok := dst.(Parser); ok {
		if err := parser.Parse(src); err != nil {
			return ParseError{Inner: err}
		}

		return nil
	}

	// dst must be pointer to a settable value
	v := reflect.ValueOf(dst)
	if v.Kind() != reflect.Pointer || v.IsNil() {
		return InvalidParseError{Type: v.Type()}
	}

	return parse(v.Elem(), src, "")
}

func parse(dst reflect.Value, src any, dstPath string) error {
	// If src is nil, we stop (do not set anything).
	if src == nil {
		return nil
	}

	if dst.CanAddr() {
		if parser, ok := dst.Addr().Interface().(Parser); ok {
			// Let the target type parse "src" however it likes
			if err := parser.Parse(src); err != nil {
				return ParseError{Inner: err, path: dstPath}
			}

			return nil
		}
	}

	vSrc := reflect.ValueOf(src)

	switch dst.Kind() {
	case reflect.Struct:
		// If dst is a struct, then src should be either a struct or a map.
		switch vSrc.Kind() {
		case reflect.Map:
			// For each key in the map, look for a field of the same name in dst.
			for _, mk := range vSrc.MapKeys() {
				// We only handle string keys here.
				keyStr, ok := mk.Interface().(string)
				if !ok {
					return ParseError{
						Msg:  fmt.Sprintf("map key %v is not a string, cannot set struct field", mk),
						path: dstPath,
					}
				}

				field := dst.FieldByName(keyStr)

				// If not found or not settable, ignore.
				if !field.IsValid() || !field.CanSet() {
					continue
				}

				if err := parse(field, vSrc.MapIndex(mk).Interface(), dstPath+"."+keyStr); err != nil {
					return err
				}
			}

		case reflect.Struct:
			// We can copy fields from one struct to the other if they match by name.
			srcType := vSrc.Type()

			for i := 0; i < srcType.NumField(); i++ {
				fieldName := srcType.Field(i).Name
				fieldDst := dst.FieldByName(fieldName)
				if !fieldDst.IsValid() || !fieldDst.CanSet() {
					continue
				}

				if err := parse(fieldDst, vSrc.Field(i).Interface(), dstPath+"."+fieldName); err != nil {
					return err
				}
			}

		default:
			return ParseError{
				Msg:  fmt.Sprintf("cannot set struct from %T", src),
				path: dstPath,
			}
		}

	case reflect.Slice:
		// If dst is a slice, src must be a slice too.
		if vSrc.Kind() != reflect.Slice {
			return ParseError{
				Msg:  fmt.Sprintf("cannot set slice from %T", src),
				path: dstPath,
			}
		}

		// Create a new slice of the appropriate type/length.
		slice := reflect.MakeSlice(dst.Type(), vSrc.Len(), vSrc.Len())

		for i := 0; i < vSrc.Len(); i++ {
			if err := parse(slice.Index(i), vSrc.Index(i).Interface(), "["+strconv.Itoa(i)+"]"); err != nil {
				return err
			}
		}

		dst.Set(slice)

	// You could handle arrays here if needed.

	default:
		// For basic types, try direct conversion.
		if vSrc.CanConvert(dst.Type()) {
			dst.Set(vSrc.Convert(dst.Type()))

			return nil
		}

		// Special-case []byte -> string
		if dst.Kind() == reflect.String &&
			vSrc.Kind() == reflect.Slice &&
			vSrc.Type().Elem().Kind() == reflect.Uint8 {
			dst.SetString(string(vSrc.Bytes()))
			return nil
		}

		return UnconvertableTypeError{
			Target:   dst.Type().String(),
			Original: reflect.TypeOf(src).String(),
		}
	}

	return nil
}
