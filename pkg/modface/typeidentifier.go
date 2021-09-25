package modface

import (
	"fmt"
	"go/ast"
	"reflect"
)

// TypeIdentifier is an interface to return both a human-readable representation of a type
// and a certain version of the type with full package name qualifiers.
type TypeIdentifier interface {
	// String returns a human-readable representation of a type.
	String() string

	// TypeID returns a certain, unmistakable type with full package name qualifiers.
	// TODO: consider making this not a part of the public interface
	typeID() string
}

type errExprTypeNotSupported struct {
	x ast.Expr
}

func (e errExprTypeNotSupported) Error() string {
	return fmt.Sprintf("type not yet supported: %s", reflect.TypeOf(e.x))
}

func parseExprToTypeIdentifier(expr ast.Expr, cache *importCache) (TypeIdentifier, error) {
	switch x := expr.(type) {
	case *ast.Ident:
		// example: type MyInt int
		// example: type MyType MyLocalType
		return LocalTypeIdentifier{TypeName: x.Name}, nil
	case *ast.SelectorExpr:
		// example: type MyContext context.Context
		return parseSelectorExprToTypeIdentifier(x, cache)
	case *ast.StarExpr:
		// example: type MyPtr *int
		// example: type MyContextPtr *context.Context
		typeID, err := parseExprToTypeIdentifier(x.X, cache)
		if err != nil {
			return nil, err
		}
		return PointerTypeIdentifier{TypeIdentifier: typeID}, nil
	case *ast.FuncType:
		// example: type MyFunc func(string) (int, error)
		return parseFuncTypeToTypeIdentifier(x, cache)
	case *ast.Ellipsis:
		typeID, err := parseExprToTypeIdentifier(x.Elt, cache)
		if err != nil {
			return nil, err
		}
		return EllipsisTypeIdentifier{TypeIdentifier: typeID}, nil
	default:
		return nil, errExprTypeNotSupported{x: x}
		// case *ast.StructType:
		// 	// example: type Duh struct{}
		// 	fallthrough
		// case *ast.InterfaceType:
		// 	// example: type Duh interface
		// 	fallthrough
		// case *ast.MapType:
		// 	// example: type Duh map[string]int
		// 	fallthrough
		// case *ast.ArrayType:
		// 	// example: type Duh []int
		// 	fallthrough
		// case *ast.ChanType:
		// 	// example: type Duh chan
		// 	fallthrough
	}
}
