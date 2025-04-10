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
		if s.Dynamic {
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

	var err error

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

			err = genValidateBody(g, path.Join(PathSegment{Name: receiver}), named.Underlying(), true)

			g.Return().Nil()
		})

	if err != nil {
		return err
	}

	return nil
}

func genValidateBody(g *jen.Group, path Path, t types.Type, addressable bool) error {
	switch t := t.(type) {
	default:
		genValidateBot(g, path, false, addressable)
		return nil

	case *types.Pointer:
		g.If(jen.Id(path.String()).Op("!=").Nil()).BlockFunc(func(g *jen.Group) {
			genValidateBot(g, path, true, addressable)
		})
		return nil

	case *types.Slice:
		return genValidateBodySlice(g, path, t, addressable)

	case *types.Map:
		return genValidateBodyMap(g, path, t)

	case *types.Struct:
		return genValidateBodyStruct(g, path, t, addressable)

	case *types.Basic:
		return nil
	}
}

func genValidateBodySlice(g *jen.Group, path Path, s *types.Slice, addressable bool) error {
	i := unique("i")

	g.For(jen.Id(i).Op(":=").Range().Id(path.String())).BlockFunc(func(g *jen.Group) {
		itemPath := path.Join(PathSegment{Name: i, Dynamic: true, Index: true})

		genValidateBody(g, itemPath, s.Elem(), addressable)
	})

	return nil
}

func genValidateBodyMap(g *jen.Group, path Path, s *types.Map) error {
	k := unique("k")

	g.For(jen.Id(k).Op(":=").Range().Id(path.String())).BlockFunc(func(g *jen.Group) {
		valuePath := path.Join(PathSegment{Name: k, Dynamic: true, Index: true})

		genValidateBody(g, valuePath, s.Elem(), false)
	})

	return nil
}

func genValidateBodyStruct(g *jen.Group, path Path, s *types.Struct, addressable bool) error {
	for field := range s.Fields() {
		if !field.Exported() {
			continue
		}
		fieldPath := path.Join(PathSegment{Name: field.Name()})

		if err := genValidateBody(g, fieldPath, field.Type(), addressable); err != nil {
			return err
		}
	}

	return nil
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
