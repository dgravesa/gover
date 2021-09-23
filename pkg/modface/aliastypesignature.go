package modface

import (
	"fmt"
	"go/ast"
	"reflect"
)

// AliasTypeSignature defines a alias type signature.
// Some examples include:
//		type MyInt          int
//		type MyIntPtr       *int
//		type MyContext      context.Context
//		type MyMathFloatPtr *math.Float64

type AliasTypeSignature struct {
	TypeName       string
	TypeIdentifier TypeIdentifier // TODO: this could be a func, or a map of a map, etc
}

func (ats AliasTypeSignature) ID() string {
	return ats.TypeName
}

func (ats AliasTypeSignature) String() string {
	return fmt.Sprintf("type %s %s", ats.TypeName, ats.TypeIdentifier)
}

func (ats AliasTypeSignature) compareString() string {
	return fmt.Sprintf("%s %s", ats.TypeName, ats.TypeIdentifier.TypeID())
}

func parseIdentTypeSignature(typeName string, v *ast.Ident) AliasTypeSignature {
	return AliasTypeSignature{
		TypeName: typeName,
		TypeIdentifier: LocalTypeIdentifier{
			TypeName:  typeName,
			IsPointer: false,
		},
	}
}

func parseSelectorTypeSignature(typeName string, s *ast.SelectorExpr, isPointer bool, cache *importCache) (AliasTypeSignature, error) {
	switch x := s.X.(type) {
	case *ast.Ident:
		pkgShortName := x.Name
		pkgFullName, found := cache.FindShortName(pkgShortName)
		if !found {
			return AliasTypeSignature{}, fmt.Errorf("import not found: %s", pkgShortName)
		}
		return AliasTypeSignature{
			TypeName: typeName,
			TypeIdentifier: ImportedTypeIdentifier{
				PackageName: pkgFullName,
				TypeName:    s.Sel.Name,
				IsPointer:   isPointer,
			},
		}, nil
	default:
		return AliasTypeSignature{}, fmt.Errorf("unsupported expr type: %s", reflect.TypeOf(s.X))
	}
}
