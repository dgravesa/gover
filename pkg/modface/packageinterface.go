package modface

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// PackageInterface represents all exports of a package.
type PackageInterface map[string]Export

func parseDir(inout ModuleInterface, basedir string, pkgdir, modname string) error {
	fset := token.NewFileSet()
	dir := filepath.Join(basedir, pkgdir)
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return err
	}

	hasExports := func(pkg *ast.Package) bool {
		if strings.HasSuffix(pkg.Name, "_test") {
			return false
		} else if !ast.PackageExports(pkg) {
			return false
		}
		return true
	}

	// parse packages
	for _, pkg := range pkgs {
		if hasExports(pkg) {
			pkgfullpath := filepath.Join(modname, pkgdir)
			pf, ok := inout[pkgfullpath]
			if !ok {
				pf = make(PackageInterface)
				inout[pkgfullpath] = pf
			}

			for _, file := range pkg.Files {
				// initialize import cache for file
				cache := newImportCache()
				for _, imp := range file.Imports {
					cache.Insert(imp)
				}

				for _, decl := range file.Decls {
					switch d := decl.(type) {
					case *ast.FuncDecl:
						fs, err := parseFuncSignature(d, cache)
						switch err.(type) {
						case nil:
							if fs.IsExported() {
								pf[fs.ID()] = fs
							}
						case errExprTypeNotSupported:
							continue // TODO: remove once all necessary types are supported
						default:
							return err
						}
					case *ast.GenDecl:
						for _, spec := range d.Specs {
							switch s := spec.(type) {
							case *ast.TypeSpec:
								ts, err := parseTypeSignature(s, cache)
								switch err.(type) {
								case nil:
									pf[ts.ID()] = ts
								case errTypeNotExported:
									continue
								case errExprTypeNotSupported:
									// TODO: just printing warning for now, remove later
									fmt.Printf("warning: %s\n", err)
								default:
									return err
								}
							case *ast.ValueSpec:
								switch d.Tok {
								case token.CONST:
									// TODO: handle package constants
									fmt.Println("consty boi") // TODO: remove
								case token.VAR:
									// TODO: handle package vars
									fmt.Println("var boi") // TODO: remove
								}
							}
						}
					}
				}
			}
		}
	}

	return nil
}
