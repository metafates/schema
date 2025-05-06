// Package parse provides parsing functionality for converting similar looking values into others with validation.
package parse

import (
	"fmt"
	"reflect"
	"strconv"

	"github.com/metafates/schema/validate"
)

type Parser interface {
	Parse(v any) error
}

// Parse attempts to copy data from src into dst. If dst implements the [Parser] interface,
// Parse simply calls dst.Parse(src). Otherwise, it uses reflection to assign fields or
// elements to dst. To succeed, dst must be a non-nil pointer to a settable value.
//
// The function supports struct-to-struct, map-to-struct, and slice-to-slice copying,
// as well as direct conversions between basic types (including []byte to string).
// If src is nil, no assignment is performed. If dst is not a valid pointer, an [InvalidParseError]
// is returned. If a type conversion is not possible, an [UnconvertableTypeError] is returned.
//
// Successfully parsed value is already validated and can be used safely.
//
// Any errors encountered during parsing are wrapped in a [ParseError].
//
// Parse also accepts options. See [Option].
func Parse(src, dst any, options ...Option) error {
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

	cfg := defaultConfig()

	for _, apply := range options {
		apply(&cfg)
	}

	if err := parse(src, v.Elem(), "", &cfg); err != nil {
		return err
	}

	if err := validate.Validate(dst); err != nil {
		return err
	}

	return nil
}

func parse(src any, dst reflect.Value, dstPath string, cfg *config) error {
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

	if vSrc.Kind() == reflect.Pointer {
		if vSrc.IsNil() {
			return nil
		}

		vSrc = vSrc.Elem()
	}

	switch dst.Kind() {
	case reflect.Struct:
		return parseToStruct(vSrc, dst, dstPath, cfg)

	case reflect.Slice:
		return parseToSlice(vSrc, dst, dstPath, cfg)

	default:
		return parseToBasic(vSrc, dst)
	}
}

func parseToBasic(src reflect.Value, dst reflect.Value) error {
	// For basic types, try direct conversion.
	if src.CanConvert(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))

		return nil
	}

	// Special-case []byte -> string
	if dst.Kind() == reflect.String &&
		src.Kind() == reflect.Slice &&
		src.Type().Elem().Kind() == reflect.Uint8 {
		dst.SetString(string(src.Bytes()))

		return nil
	}

	return ParseError{
		Inner: UnconvertableTypeError{
			Target:   dst.Type().String(),
			Original: reflect.TypeOf(src).String(),
		},
	}
}

func parseToStruct(src reflect.Value, dst reflect.Value, dstPath string, cfg *config) error {
	// If dst is a struct, then src should be either a struct or a map.
	switch src.Kind() {
	case reflect.Map:
		return parseMapToStruct(src, dst, dstPath, cfg)

	case reflect.Struct:
		return parseStructToStruct(src, dst, dstPath, cfg)

	default:
		return ParseError{
			Msg:  fmt.Sprintf("cannot set struct from %T", src),
			path: dstPath,
		}
	}
}

func parseStructToStruct(src reflect.Value, dst reflect.Value, dstPath string, cfg *config) error {
	// We can copy fields from one struct to the other if they match by name.
	srcType := src.Type()

	for i := range srcType.NumField() {
		fieldSrc := srcType.Field(i)
		if !fieldSrc.IsExported() {
			continue
		}

		fieldName := cfg.RenameFunc(fieldSrc.Name)

		fieldDst := dst.FieldByName(fieldName)

		if !fieldDst.IsValid() || !fieldDst.CanSet() {
			if cfg.DisallowUnknownFields {
				return ParseError{
					Inner: UnknownFieldError{Name: fieldName},
				}
			}

			continue
		}

		if err := parse(src.Field(i).Interface(), fieldDst, dstPath+"."+fieldName, cfg); err != nil {
			return err
		}
	}

	return nil
}

func parseMapToStruct(src reflect.Value, dst reflect.Value, dstPath string, cfg *config) error {
	// For each key in the map, look for a field of the same name in dst.
	for _, mk := range src.MapKeys() {
		// We only handle string keys here.
		keyStr, ok := mk.Interface().(string)
		if !ok {
			return ParseError{
				Msg:  fmt.Sprintf("map key %v is not a string, cannot set struct field", mk),
				path: dstPath,
			}
		}

		keyStr = cfg.RenameFunc(keyStr)

		field := dst.FieldByName(keyStr)

		// If not found or not settable, ignore.
		if !field.IsValid() || !field.CanSet() {
			if cfg.DisallowUnknownFields {
				return ParseError{
					Inner: UnknownFieldError{Name: keyStr},
				}
			}

			continue
		}

		if err := parse(src.MapIndex(mk).Interface(), field, dstPath+"."+keyStr, cfg); err != nil {
			return err
		}
	}

	return nil
}

func parseToSlice(src reflect.Value, dst reflect.Value, dstPath string, cfg *config) error {
	// If dst is a slice, src must be a slice too.
	if src.Kind() != reflect.Slice {
		return ParseError{
			Msg:  fmt.Sprintf("cannot set slice from %T", src),
			path: dstPath,
		}
	}

	// Create a new slice of the appropriate type/length.
	slice := reflect.MakeSlice(dst.Type(), src.Len(), src.Len())

	for i := range src.Len() {
		if err := parse(src.Index(i).Interface(), slice.Index(i), "["+strconv.Itoa(i)+"]", cfg); err != nil {
			return err
		}
	}

	dst.Set(slice)

	return nil
}
