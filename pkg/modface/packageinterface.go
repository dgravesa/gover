package modface

import (
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
				for _, decl := range file.Decls {
					switch v := decl.(type) {
					case *ast.FuncDecl:
						fs := ParseFuncSignature(v)
						funcExported := ast.IsExported(fs.Name)
						recvNotAnonymous := !fs.Receiver.IsDefined() || fs.Receiver.IsExported()
						if funcExported && recvNotAnonymous {
							pf[fs.ID()] = fs
						}
					}
				}
			}
		}
	}

	return nil
}
