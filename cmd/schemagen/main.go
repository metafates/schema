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

var typeNames = flag.String("type", "", "comma-separated list of type names; must be set")

func Usage() {
	printf := func(format string, a ...any) {
		fmt.Fprintf(os.Stderr, format, a...)
	}

	printf("Schemagen is a tool to generate Go code for field-traversal validation\n")
	printf("Usage of %s:\n", os.Args[0])
	printf("\tschemagen [flags] -type T [directory]\n")
	printf("For more information, see:\n")
	printf("\thttps://github.com/metafates/schema\n")
	printf("Flags:\n")

	flag.PrintDefaults()
}

func main() {
	log.SetFlags(0)
	log.SetPrefix("schemagen: ")

	flag.Usage = Usage
	flag.Parse()

	if len(*typeNames) == 0 {
		flag.Usage()
		os.Exit(2)
	}

	args := flag.Args()
	if len(args) == 0 {
		// Default: process whole package in current directory.
		args = []string{"."}
	}

	for _, pattern := range args {
		genPackage(pattern)
	}
}

func genPackage(pattern string) {
	pkg := parsePackage(pattern)

	scope := pkg.Types.Scope()

	f := jen.NewFilePathName(pkg.PkgPath, pkg.Name)
	f.HeaderComment(fmt.Sprintf(`Code generated by "schemagen %s". DO NOT EDIT.`, strings.Join(os.Args[1:], " ")))

	for _, name := range strings.Split(*typeNames, ",") {
		obj := scope.Lookup(name)
		if obj == nil {
			log.Fatalf("unknown type: %s\n", name)
		}

		genType(f, obj.Type().(*types.Named))
	}

	if err := f.Save(outputPath(pkg.Dir)); err != nil {
		log.Fatalln(err)
	}
}

func outputPath(dir string) string {
	return filepath.Join(dir, "schemagen.go")
}

func parsePackage(pattern string) *packages.Package {
	cfg := packages.Config{
		Mode:  packages.LoadAllSyntax,
		Tests: false,
	}

	pkgs, err := packages.Load(&cfg, pattern)
	if err != nil {
		log.Fatalln(err)
	}

	if len(pkgs) != 1 {
		log.Fatalf("error: %d packages found", len(pkgs))
	}

	// if packages.PrintErrors(pkgs) > 0 {
	// 	os.Exit(1)
	// }

	return pkgs[0]
}

func genType(f *jen.File, named *types.Named) {
	genLock(f, named)
	genValidate(f, named)
}
