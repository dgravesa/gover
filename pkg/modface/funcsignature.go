package modface

import (
	"fmt"
	"go/ast"
	"strings"
)

// FuncSignature defines a function or method export.
type FuncSignature struct {
	// Name is the name of the function or method.
	Name string

	// Receiver is nil if export is a function, or either a LocalTypeIdentifier or a
	// PointerTypeIdentifier wrapping a LocalTypeIdentifier if export is a method.
	Receiver TypeIdentifier

	// Type defines the parameter and result types of the function.
	Type *FuncTypeIdentifier
}

// ID returns a unique identifier for the function signature.
// A package should only have one function with this particular ID.
func (fs FuncSignature) ID() string {
	switch r := fs.Receiver.(type) {
	case nil:
		return fs.Name
	case *PointerTypeIdentifier:
		return fmt.Sprintf("%s.%s", r.TypeIdentifier, fs.Name)
	default:
		return fmt.Sprintf("%s.%s", r, fs.Name)
	}
}

func (fs FuncSignature) String() string {
	var sb strings.Builder

	sb.WriteString("func ")

	// if method, add the receiver to the signature
	if fs.Receiver != nil {
		sb.WriteString(fmt.Sprintf("(%s) ", fs.Receiver))
	}

	// write function name and base type string
	sb.WriteString(fmt.Sprintf("%s%s", fs.Name, fs.Type.formBaseTypeString()))

	return sb.String()
}

// IsExported returns true if the function or method is exported, otherwise false.
func (fs FuncSignature) IsExported() bool {
	return ast.IsExported(fs.Name) && (fs.Receiver == nil || ast.IsExported(fs.Receiver.String()))
}

func (fs FuncSignature) compareString() string {
	// TODO: come up with cleaner standard for compare strings across different types of exports.
	if fs.Receiver != nil {
		return fmt.Sprintf("(%s) %s %s", fs.Receiver, fs.Name, fs.Type.typeID())
	}
	return fmt.Sprintf("%s %s", fs.Name, fs.Type.typeID())
}

func parseFuncSignature(decl *ast.FuncDecl, cache *importCache) (*FuncSignature, error) {
	var fs FuncSignature
	var err error

	fs.Name = decl.Name.Name

	if decl.Recv != nil && len(decl.Recv.List) > 0 {
		fs.Receiver, err = parseExprToTypeIdentifier(decl.Recv.List[0].Type, cache)
		if err != nil {
			return nil, err
		}
	}

	fs.Type, err = parseFuncTypeToFuncTypeIdentifier(decl.Type, cache)
	if err != nil {
		return nil, err
	}

	return &fs, nil
}
