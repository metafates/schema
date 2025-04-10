package main

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/metafates/schema/internal/typeconv"
)

func genLock(f *jen.File, named *types.Named) {
	switch t := named.Underlying().(type) {
	default:
	case *types.Struct, *types.Map, *types.Slice:
		f.Comment("Ensure types are not changed")
		f.Func().Id("_").Params().BlockFunc(func(g *jen.Group) {
			converter := typeconv.NewTypeConverter()
			typeCode := converter.ConvertType(t)

			converter.AddImports(f)

			g.Type().Id("locked").Add(typeCode)

			g.Comment("Compiler error signifies that the type definition have changed.")
			g.Comment("Re-run the schemagen command to regenerate this file.")
			g.Id("_").Op("=").Id("locked").Params(jen.Id(named.Obj().Name()).Values())
		})
	}
}
