package main

import (
	"go/types"

	"github.com/dave/jennifer/jen"
	"github.com/metafates/schema/internal/typeconv"
)

func genLock(f *jen.File, named *types.Named) {
	f.Comment("Ensure types are not changed")
	f.Func().Id("_").Params().BlockFunc(func(g *jen.Group) {
		underlying := named.Underlying()

		converter := typeconv.NewTypeConverter()
		typeCode := converter.ConvertType(underlying)

		converter.AddImports(f)

		g.Type().Id("locked").Add(typeCode)

		varName := "v"

		g.Var().Id(varName).Id(named.Obj().Name())

		g.Comment("Compiler error signifies that the type definition have changed.")
		g.Comment("Re-run the schemagen command to regenerate this file.")
		g.Id("_").Op("=").Id("locked").Params(jen.Id(varName))
	})
}
