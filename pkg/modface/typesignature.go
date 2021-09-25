package modface

import (
	"fmt"
	"go/ast"
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
	return fmt.Sprintf("%s %s", ts.Name, ts.TypeIdentifier.typeID())
}

type errTypeNotExported struct {
	typeName string
}

func (e errTypeNotExported) Error() string {
	return fmt.Sprintf("type not exported: %s", e.typeName)
}

func parseTypeSignature(spec *ast.TypeSpec, cache *importCache) (*TypeSignature, error) {
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
