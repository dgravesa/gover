package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/dgravesa/gover/pkg/modface"
	"github.com/spf13/cobra"
)

var printModpaths = []string{"."}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print module interface",
	Long: `
	Print all exports of a module
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			printModpaths = args
		}

		for _, modpath := range printModpaths {
			err := printModpath(modpath)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func printModpath(moddir string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, moddir, nil, 0)
	if err != nil {
		return err
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

	return nil
}
