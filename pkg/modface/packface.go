package modface

import (
	"go/ast"
	"go/parser"
	"go/token"
)

// PackFace represents all exports of a package.
type PackFace map[string]Face

func parseDir(dir string, inout ModFace) error {
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, dir, nil, 0)
	if err != nil {
		return err
	}

	// parse packages
	for _, pkg := range pkgs {

		// TODO: do this at the module level
		if ast.PackageExports(pkg) {
			pf, ok := inout[pkg.Name]
			if !ok {
				pf = make(PackFace)
				inout[pkg.Name] = pf
			}

			for _, file := range pkg.Files {
				for _, decl := range file.Decls {
					switch v := decl.(type) {
					case *ast.FuncDecl:
						fsig := ParseFuncSig(v)
						// fsigs[fsig.ID()] = fsig
						pf[fsig.ID()] = fsig
					}
				}
			}
		}
	}

	return nil
}
