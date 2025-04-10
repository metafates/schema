package main

import (
	"fmt"
	"go/types"

	"github.com/dave/jennifer/jen"
)

const validatePkg = "github.com/metafates/schema/validate"

func genValidate(f *jen.File, named *types.Named) {
	const receiver = "x"

	receiverPtr := "*"

	switch named.Underlying().(type) {
	case *types.Slice, *types.Map:
		receiverPtr = ""
	}

	f.Comment("Validate implementes [validate.Validateable]")
	f.
		Func().
		Params(jen.Id(receiver).Op(receiverPtr).Id(named.Obj().Name())).
		Id("Validate").
		Params().
		Error().
		BlockFunc(func(g *jen.Group) {
			gen := validateGenerator{counter: make(map[string]int)}

			gen.gen(
				g,
				Path{}.Join(PathSegment{Name: receiver}),
				named.Underlying(),
				false,
				true,
			)

			g.Return().Nil()
		})
}

type validateGenerator struct {
	counter map[string]int
}

func (vg *validateGenerator) gen(g *jen.Group, path Path, t types.Type, isPtr, addressable bool) bool {
	switch t := t.(type) {
	default:
		vg.genBot(g, path, isPtr, addressable)
		return true

	case *types.Pointer:
		return vg.genPointer(g, path, t, addressable)

	case *types.Slice:
		return vg.genSlice(g, path, t, isPtr, addressable)

	case *types.Map:
		return vg.genMap(g, path, t, isPtr)

	case *types.Struct:
		return vg.genStruct(g, path, t, isPtr, addressable)

	case *types.Basic:
		return false
	}
}

func (vg *validateGenerator) genPointer(g *jen.Group, path Path, s *types.Pointer, addressable bool) bool {
	var generated bool

	ifBody := jen.BlockFunc(func(g *jen.Group) {
		generated = vg.gen(g, path, s.Elem(), true, addressable)
	})

	if generated {
		g.If(jen.Id(path.String()).Op("!=").Nil()).Block(ifBody)
	}

	return generated
}

func (vg *validateGenerator) genSlice(g *jen.Group, path Path, s *types.Slice, isPtr, addressable bool) bool {
	i := vg.unique("i")

	var generatedAny bool

	loopBody := jen.BlockFunc(func(g *jen.Group) {
		itemPath := path.Join(PathSegment{Name: i, Dynamic: true, Index: true})

		if vg.gen(g, itemPath, s.Elem(), isPtr, addressable) {
			generatedAny = true
		}
	})

	if generatedAny {
		g.For(jen.Id(i).Op(":=").Range().Id(path.String())).Block(loopBody)
	}

	return generatedAny
}

func (vg *validateGenerator) genMap(g *jen.Group, path Path, s *types.Map, isPtr bool) bool {
	k := vg.unique("k")

	var generatedAny bool

	loopBody := jen.BlockFunc(func(g *jen.Group) {
		valuePath := path.Join(PathSegment{Name: k, Dynamic: true, Index: true})

		if vg.gen(g, valuePath, s.Elem(), isPtr, false) {
			generatedAny = true
		}
	})

	if generatedAny {
		g.For(jen.Id(k).Op(":=").Range().Id(path.String())).Block(loopBody)
	}

	return generatedAny
}

func (vg *validateGenerator) genStruct(g *jen.Group, path Path, s *types.Struct, isPtr, addressable bool) bool {
	var generatedAny bool

	for field := range s.Fields() {
		if !field.Exported() {
			continue
		}
		fieldPath := path.Join(PathSegment{Name: field.Name()})

		if vg.gen(g, fieldPath, field.Type(), isPtr, addressable) {
			generatedAny = true
		}
	}

	return generatedAny
}

func (vg *validateGenerator) genBot(g *jen.Group, path Path, isPtr, addressable bool) {
	errName := vg.unique("err")

	var value string

	if isPtr || !addressable {
		value = path.String()
	} else {
		value = "&" + path.String()
	}

	valueName := vg.unique("v")

	g.Id(valueName).Op(":=").Id(value)

	g.Id(errName).Op(":=").Qual(validatePkg, "Validate").Call(jen.Id(valueName))

	if !addressable {
		g.Id(path.String()).Op("=").Id(valueName)
	}

	g.If(jen.Id(errName).Op("!=").Nil()).Block(
		jen.Return().
			Qual(validatePkg, "ValidationError").
			Values(jen.Dict{
				jen.Id("Inner"): jen.Id(errName),
			}).
			Dot("WithPath").
			Call(jen.Qual("fmt", "Sprintf").CallFunc(func(g *jen.Group) {
				format, args := path.printf()

				g.Lit(format)

				for _, a := range args {
					g.Id(a)
				}
			})),
	)
}

func (vg *validateGenerator) unique(id string) string {
	// we could use a single counter for all ids
	// but using counter for each id generates better looking code
	count, ok := vg.counter[id]
	if !ok {
		vg.counter[id] = count
	}

	vg.counter[id]++

	return fmt.Sprintf("%s%d", id, count)
}
