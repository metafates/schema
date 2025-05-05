package main

import (
	"slices"
	"strings"
)

type PathSegment struct {
	Name    string
	Index   bool
	Dynamic bool
}

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
