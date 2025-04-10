package main

import (
	"fmt"
	"go/types"
	"slices"
	"strings"
	"sync/atomic"

	"github.com/dave/jennifer/jen"
)

var counter atomic.Int32

const validatePkg = "github.com/metafates/schema/validate"

type Path struct {
	Segments []PathSegment
}

func (p Path) printf() (format string, args []string) {
	var formatBuilder strings.Builder

	for _, s := range p.Segments[1:] {
		var prefix, suffix string

		if s.Index {
			prefix = "["
			suffix = "]"
		} else {
			prefix = "."
		}

		formatBuilder.WriteString(prefix)

		if s.Dynamic {
			formatBuilder.WriteString("%v")
			args = append(args, s.Name)
		} else {
			formatBuilder.WriteString(s.Name)
		}

		formatBuilder.WriteString(suffix)
	}

	return formatBuilder.String(), args
}

func (p Path) Join(segment PathSegment) Path {
	return Path{
		Segments: append(slices.Clone(p.Segments), segment),
	}
}

func (p Path) String() string {
	if len(p.Segments) == 0 {
		return ""
	}

	root := p.Segments[0].Name

	var rest strings.Builder

	rest.Grow(50)

	for _, s := range p.Segments[1:] {
		if s.Index {
			rest.WriteString("[" + s.Name + "]")
		} else {
			rest.WriteString("." + s.Name)
		}
	}

	return root + rest.String()
}

type PathSegment struct {
	Name    string
	Index   bool
	Dynamic bool
}

func genValidate(f *jen.File, named *types.Named) error {
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
			var path Path

			genValidateBody(g, path.Join(PathSegment{Name: receiver}), named.Underlying(), false, true)

			g.Return().Nil()
		})

	return nil
}

func genValidateBody(g *jen.Group, path Path, t types.Type, isPtr, addressable bool) bool {
	switch t := t.(type) {
	default:
		genValidateBot(g, path, isPtr, addressable)
		return true

	case *types.Pointer:
		return genValidateBodyPointer(g, path, t, addressable)

	case *types.Slice:
		return genValidateBodySlice(g, path, t, isPtr, addressable)

	case *types.Map:
		return genValidateBodyMap(g, path, t, isPtr)

	case *types.Struct:
		return genValidateBodyStruct(g, path, t, isPtr, addressable)

	case *types.Basic:
		return false
	}
}

func genValidateBodyPointer(g *jen.Group, path Path, s *types.Pointer, addressable bool) bool {
	var generated bool

	ifBody := jen.BlockFunc(func(g *jen.Group) {
		generated = genValidateBody(g, path, s.Elem(), true, addressable)
	})

	if generated {
		g.If(jen.Id(path.String()).Op("!=").Nil()).Block(ifBody)
	}

	return generated
}

func genValidateBodySlice(g *jen.Group, path Path, s *types.Slice, isPtr, addressable bool) bool {
	i := unique("i")

	var generatedAny bool

	loopBody := jen.BlockFunc(func(g *jen.Group) {
		itemPath := path.Join(PathSegment{Name: i, Dynamic: true, Index: true})

		if genValidateBody(g, itemPath, s.Elem(), isPtr, addressable) {
			generatedAny = true
		}
	})

	if generatedAny {
		g.For(jen.Id(i).Op(":=").Range().Id(path.String())).Block(loopBody)
	}

	return generatedAny
}

func genValidateBodyMap(g *jen.Group, path Path, s *types.Map, isPtr bool) bool {
	k := unique("k")

	var generatedAny bool

	loopBody := jen.BlockFunc(func(g *jen.Group) {
		valuePath := path.Join(PathSegment{Name: k, Dynamic: true, Index: true})

		if genValidateBody(g, valuePath, s.Elem(), isPtr, false) {
			generatedAny = true
		}
	})

	if generatedAny {
		g.For(jen.Id(k).Op(":=").Range().Id(path.String())).Block(loopBody)
	}

	return generatedAny
}

func genValidateBodyStruct(g *jen.Group, path Path, s *types.Struct, isPtr, addressable bool) bool {
	var generatedAny bool

	for field := range s.Fields() {
		if !field.Exported() {
			continue
		}
		fieldPath := path.Join(PathSegment{Name: field.Name()})

		if genValidateBody(g, fieldPath, field.Type(), isPtr, addressable) {
			generatedAny = true
		}
	}

	return generatedAny
}

func genValidateBot(g *jen.Group, path Path, isPtr, addressable bool) {
	errName := unique("err")

	var value string

	if isPtr || !addressable {
		value = path.String()
	} else {
		value = "&" + path.String()
	}

	valueName := unique("v")

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

func unique(id string) string {
	return fmt.Sprintf("%s%d", id, counter.Add(1))
}
