package modface

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"
)

// PackFace represents all exports of a package.
type PackFace map[string]Face

func parseDir(inout ModFace, dir string, modname string) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return err
	}

	isFacePackage := func(pkg *ast.Package) bool {
		if strings.HasSuffix(pkg.Name, "_test") {
			return false
		} else if !ast.PackageExports(pkg) {
			return false
		}
		return true
	}

	// parse packages
	for _, pkg := range pkgs {
		if isFacePackage(pkg) {
			pkgfullpath := filepath.Join(modname, dir)
			pf, ok := inout[pkgfullpath]
			if !ok {
				pf = make(PackFace)
				inout[pkgfullpath] = pf
			}

			for _, file := range pkg.Files {
				for _, decl := range file.Decls {
					switch v := decl.(type) {
					case *ast.FuncDecl:
						fsig := ParseFuncSig(v)
						pf[fsig.ID()] = fsig
					}
				}
			}
		}
	}

	return nil
}
