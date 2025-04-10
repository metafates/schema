package main

import (
	"flag"
	"fmt"
	"go/types"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dave/jennifer/jen"
	"golang.org/x/tools/go/packages"
)

type Env struct {
	File string
}

func NewEnv() Env {
	return Env{
		File: os.Getenv("GOFILE"),
	}
}

// Output returns output path for the generated file
func (e Env) Output() string {
	dir := filepath.Dir(e.File)
	ext := filepath.Ext(e.File)
	stem := strings.TrimSuffix(filepath.Base(e.File), ext)

	return filepath.Join(dir, stem+".schema"+ext)
}

type arrayFlag []string

func (f *arrayFlag) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *arrayFlag) Set(value string) error {
	for _, subvalue := range strings.Split(value, ",") {
		*f = append(*f, subvalue)
	}

	return nil
}

func main() {
	var typeNames arrayFlag

	flag.Var(&typeNames, "type", "comma-separated list of type names; must be set")

	flag.Parse()

	if len(typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	env := NewEnv()

	path := "."

	pkgs, err := packages.Load(
		&packages.Config{
			Mode: packages.LoadAllSyntax,
		},
		path,
	)
	if err != nil {
		log.Fatalln(err)
	}

	if packages.PrintErrors(pkgs) > 0 {
		os.Exit(1)
	}

	pkg := pkgs[0]

	scope := pkg.Types.Scope()

	f := jen.NewFilePathName(pkg.PkgPath, pkg.Name)
	f.HeaderComment(fmt.Sprintf(`Code generated by "schemagen %s". DO NOT EDIT.`, strings.Join(os.Args[1:], " ")))

	for _, name := range typeNames {
		obj := scope.Lookup(name)
		if obj == nil {
			log.Fatalf("unknown type: %s\n", name)
		}

		err = handleObject(f, obj.Type().(*types.Named))
		if err != nil {
			log.Fatalln(err)
		}
	}

	err = f.Save(env.Output())
	if err != nil {
		log.Fatalln(err)
	}
}

func handleObject(f *jen.File, named *types.Named) error {
	return genValidate(f, named)
}
