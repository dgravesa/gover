package modface

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"reflect"
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
func (iti ImportedTypeIdentifier) TypeID() string {
	return fmt.Sprintf("%s.%s", iti.PackageName, iti.TypeName)
}

func parseSelectorExprToImportedTypeIdentifier(s *ast.SelectorExpr, cache *importCache) (ImportedTypeIdentifier, error) {
	switch x := s.X.(type) {
	case *ast.Ident:
		pkgShortName := x.Name
		pkgFullName, found := cache.FindShortName(pkgShortName)
		if !found {
			return ImportedTypeIdentifier{}, fmt.Errorf("import not found: %s", pkgShortName)
		}
		return ImportedTypeIdentifier{
			PackageName: pkgFullName,
			TypeName:    s.Sel.Name,
		}, nil
	default:
		return ImportedTypeIdentifier{}, fmt.Errorf("unsupported expr type: %s", reflect.TypeOf(s.X))
	}
}
