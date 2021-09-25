package modface

import (
	"fmt"
	"go/ast"
	"path/filepath"
)

// ImportedTypeIdentifier corresponds to a type imported from another package.
type ImportedTypeIdentifier struct {
	PackageName string
	TypeName    string
}

func (iti ImportedTypeIdentifier) String() string {
	pkgShortName := filepath.Base(iti.PackageName)
	return fmt.Sprintf("%s.%s", pkgShortName, iti.TypeName)
}

// TypeID returns a unique ID for the type across any go packages.
func (iti ImportedTypeIdentifier) typeID() string {
	return fmt.Sprintf("%s.%s", iti.PackageName, iti.TypeName)
}

type errImportNotFound struct {
	pkgShortName string
}

func (e errImportNotFound) Error() string {
	return fmt.Sprintf("import not found: %s", e.pkgShortName)
}

func parseSelectorExprToTypeIdentifier(s *ast.SelectorExpr, cache *importCache) (TypeIdentifier, error) {
	switch x := s.X.(type) {
	case *ast.Ident:
		pkgShortName := x.Name
		pkgFullName, found := cache.FindShortName(pkgShortName)
		if !found {
			return nil, errImportNotFound{pkgShortName: pkgShortName}
		}
		return &ImportedTypeIdentifier{
			PackageName: pkgFullName,
			TypeName:    s.Sel.Name,
		}, nil
	default:
		return nil, errExprTypeNotSupported{x: s.X}
	}
}
