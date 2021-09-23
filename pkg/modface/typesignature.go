package modface

import (
	"fmt"
	"go/ast"
	"reflect"
)

// TODO:
// types of "types":
// type Duh int
// type Duh struct {
//		...
// }
// type Duh interface {
//		...
// }
type TypeSignature interface {
	Export
	// TODO: this may not be needed, uncomment otherwise
	// TypeIdentifier
}

// TODO: may need to incorporate an import cache into arguments.
func parseTypeSignatureAndUpdatePackageInterface(
	inout PackageInterface, spec *ast.TypeSpec, packageName string, prefix string, cache *importCache) error {
	typeName := spec.Name.Name

	// do not continue if type is not exported
	if !ast.IsExported(typeName) {
		return nil
	}

	switch v := spec.Type.(type) {
	case *ast.StructType:
		// TODO: implement
		// example: type Duh struct{}
	case *ast.InterfaceType:
		// TODO: implement
		// example: type Duh interface
	case *ast.MapType:
		// TODO: implement
		// example: type Duh map[string]int
	case *ast.ArrayType:
		// TODO: implement
		// example: type Duh []int
	case *ast.Ident:
		// example: type MyInt int
		// TODO: this may not need module cache at all
		ts, err := parseIdentTypeSignature(typeName, v, cache)
		if err != nil {
			return err
		}
		inout[ts.ID()] = ts
	case *ast.StarExpr:
		// TODO: implement
		// example: type MyPtr *int
		// example: type MyContextPtr *context.Context
	case *ast.SelectorExpr:
		// TODO: implement
		// example: type MyFloat math.Float64
	case *ast.FuncType:
		// TODO: implement
		// example: type DuhFunc func(int) int
	case *ast.ChanType:
		// TODO: implement
		// example: type Duh chan
	}

	// TODO: figure out what kind of type it is
	fmt.Printf("%s: %s\n", typeName, reflect.TypeOf(spec.Type))

	return nil
}

// TODO: implement, then remove all of these
// type MyInt int // implemented
// type MyPtr *int
// type MyFloat math.Float64
// type MyContextPtr *context.Context
