package modface

import (
	"fmt"
	"go/ast"
)

// TODO: Ident may not need BaseTypeIdentifier because it doesn't seem to ever need package detail
type IdentTypeSignature struct {
	TypeName string
	BaseType BaseTypeIdentifier
}

func (is IdentTypeSignature) ID() string {
	return is.TypeName
}

func (is IdentTypeSignature) String() string {
	return fmt.Sprintf("type %s %s", is.TypeName, is.BaseType)
}

func (is IdentTypeSignature) compareString() string {
	return fmt.Sprintf("%s %s", is.TypeName, is.BaseType.TypeID())
}

func parseIdentTypeSignature(typeName string, v *ast.Ident, cache *importCache) (IdentTypeSignature, error) {
	baseType, err := parseBaseTypeIdentifier(v, cache)
	return IdentTypeSignature{
		TypeName: typeName,
		BaseType: baseType,
	}, err
}
