package main

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"

	"github.com/dgravesa/gover/pkg/modface"
)

func main() {
	errOccurred := false
	flag.Parse()

	for _, moddir := range flag.Args() {
		fset := token.NewFileSet()
		pkgs, err := parser.ParseDir(fset, moddir, nil, 0)
		if err != nil {
			errOccurred = true
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		// parse packages
		for _, pkg := range pkgs {
			// TODO: do this at the module level
			fsigs := make(map[string]string)

			if ast.PackageExports(pkg) {
				for _, file := range pkg.Files {
					for _, decl := range file.Decls {
						switch v := decl.(type) {
						case *ast.FuncDecl:
							fsig := modface.ParseFuncSig(v)
							fsigs[fsig.ID()] = fsig.String()
						}
					}
				}
			}

			// TODO: do this at the module level
			for _, fs := range fsigs {
				fmt.Printf("%s\n", fs)
			}
		}
	}

	if errOccurred {
		os.Exit(1)
	}
	os.Exit(0)
}
