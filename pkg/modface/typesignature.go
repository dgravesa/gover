package modface

import (
	"fmt"
	"go/ast"
	"reflect"
)

// TypeSignature defines any type export.
type TypeSignature struct {
	Name           string
	TypeIdentifier TypeIdentifier
}

func (ts TypeSignature) String() string {
	return fmt.Sprintf("type %s %s", ts.Name, ts.TypeIdentifier)
}

// ID returns a unique identifier for an export.
func (ts TypeSignature) ID() string {
	return ts.Name
}

func (ts TypeSignature) compareString() string {
	return fmt.Sprintf("%s %s", ts.Name, ts.TypeIdentifier.TypeID())
}

type errTypeNotExported struct {
	typeName string
}

func (e errTypeNotExported) Error() string {
	return fmt.Sprintf("type not exported: %s", e.typeName)
}

func parseTypeSignature(spec *ast.TypeSpec, cache *importCache) (Export, error) {
	typeName := spec.Name.Name

	// do not continue if type is not exported
	if !ast.IsExported(typeName) {
		return nil, errTypeNotExported{typeName: typeName}
	}

	typeID, err := parseExprToTypeIdentifier(spec.Type, cache)
	if err != nil {
		return nil, err
	}

	return &TypeSignature{
		Name:           typeName,
		TypeIdentifier: typeID,
	}, nil
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
		return parseSelectorExprToImportedTypeIdentifier(x, cache)
	case *ast.StarExpr:
		// example: type MyPtr *int
		// example: type MyContextPtr *context.Context
		typeID, err := parseExprToTypeIdentifier(x.X, cache)
		if err != nil {
			return nil, err
		}
		return PointerTypeIdentifier{TypeIdentifier: typeID}, nil
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
		// case *ast.FuncType:
		// 	// example: type DuhFunc func(int) int
		// 	fallthrough
		// case *ast.ChanType:
		// 	// example: type Duh chan
		// 	fallthrough
	}
}
