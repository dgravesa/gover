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
						fs := parseFuncSignature(d)
						funcExported := ast.IsExported(fs.Name)
						recvNotAnonymous := !fs.Receiver.IsDefined() || fs.Receiver.IsExported()
						if funcExported && recvNotAnonymous {
							pf[fs.ID()] = fs
						}
					case *ast.GenDecl:
						for _, spec := range d.Specs {
							switch s := spec.(type) {
							case *ast.ImportSpec:
								// TODO: does this ever happen?
								cache.Insert(s)
							case *ast.TypeSpec:
								// ts := parseTypeSignature(s, pkgfullpath)
								// NOTE: because of recursive nature of parsing, it's more
								// efficient to update package interface while parsing the
								// type signature, its fields, and potential subfields.
								err := parseTypeSignatureAndUpdatePackageInterface(
									pf, s, pkgfullpath, "", cache)
								if err != nil {
									return err
								}
								// if ast.IsExported(ts.Name) {
								// 	pf[ts.ID()] = ts
								// 	insertRecursiveFieldSignatures(pf, ts.Fields, "")
								// 	for _, field := range ts.Fields {

								// 	}
								// }
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
