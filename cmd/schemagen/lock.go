package main

import (
	"go/types"

	"github.com/dave/jennifer/jen"
)

// TODO: tests are needed!

func lockType(f *jen.File, named *types.Named) {
	f.Comment("Ensure types are not changed")
	f.Func().Id("_").Params().BlockFunc(func(g *jen.Group) {
		switch t := named.Underlying().(type) {
		case *types.Struct, *types.Map, *types.Slice:
			converter := &TypeConverter{imports: make(map[string]string)}
			typeCode := converter.ConvertType(t)

			// Add imports to the file
			for path, name := range converter.imports {
				f.ImportName(path, name)
			}

			g.Type().Id("locked").Add(typeCode)

			g.Comment("Compiler error signifies that the type definition have changed.")
			g.Comment("Re-run the schemagen command to regenerate validators.")
			g.Id("_").Op("=").Id("locked").Params(jen.Id(named.Obj().Name()).Values())
		}
	})
}

// NOTE: generated with AI. Validate it

type TypeConverter struct {
	imports map[string]string
}

func (c *TypeConverter) ConvertType(t types.Type) jen.Code {
	switch t := t.(type) {
	case *types.Named:
		pkg := t.Obj().Pkg()
		if pkg != nil && pkg.Path() != "" {
			// Record the import
			c.imports[pkg.Path()] = pkg.Name()

			// Create a qualified reference
			qual := jen.Qual(pkg.Path(), t.Obj().Name())

			// Handle type arguments if present
			typeArgs := t.TypeArgs()
			if typeArgs != nil && typeArgs.Len() > 0 {
				var args []jen.Code
				for i := 0; i < typeArgs.Len(); i++ {
					args = append(args, c.ConvertType(typeArgs.At(i)))
				}
				return qual.Index(args...)
			}

			return qual
		}
		return jen.Id(t.Obj().Name())

	case *types.Basic:
		return jen.Id(t.Name())

	case *types.Pointer:
		return jen.Op("*").Add(c.ConvertType(t.Elem()))

	case *types.Slice:
		return jen.Index().Add(c.ConvertType(t.Elem()))

	case *types.Map:
		return jen.Map(c.ConvertType(t.Key())).Add(c.ConvertType(t.Elem()))

	case *types.Struct:
		fields := make([]jen.Code, 0, t.NumFields())
		for i := 0; i < t.NumFields(); i++ {
			field := t.Field(i)
			tag := t.Tag(i)
			fieldCode := jen.Id(field.Name()).Add(c.ConvertType(field.Type()))
			if tag != "" {
				tagMap := parseStructTag(tag)
				fieldCode = fieldCode.Tag(tagMap)
			}

			fields = append(fields, fieldCode)
		}

		return jen.Struct(fields...)

	default:
		// Fallback for any other types
		return jen.Id(t.String())
	}
}

// parseStructTag parses a raw struct tag string into a map[string]string
func parseStructTag(tag string) map[string]string {
	tags := make(map[string]string)

	// Simple state machine to parse tags
	for tag != "" {
		// Skip leading space
		i := 0
		for i < len(tag) && tag[i] == ' ' {
			i++
		}
		tag = tag[i:]
		if tag == "" {
			break
		}

		// Scan to colon
		i = 0
		for i < len(tag) && tag[i] != ':' {
			i++
		}
		if i >= len(tag) {
			break
		}

		name := tag[:i]
		tag = tag[i+1:]

		// Scan to closing quote, handling escaped quotes
		if tag[0] != '"' {
			break
		}

		i = 1
		for i < len(tag) {
			if tag[i] == '"' && tag[i-1] != '\\' {
				break
			}
			i++
		}
		if i >= len(tag) {
			break
		}

		value := tag[1:i]
		tag = tag[i+1:]

		tags[name] = value
	}

	return tags
}
