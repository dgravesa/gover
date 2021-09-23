package modface

import (
	"fmt"
	"go/ast"
	"strings"
)

// FuncSignature defines a function signature.
// The ID is a short signature that should be uniquely identifying.
// The Signature is a complete representation of the function's interface
// and should be directly comparable between different commits to ensure
// that backwards compatibility is maintained.
type FuncSignature struct {
	Name     string
	Receiver Type
	Params   TypeList
	Results  TypeList
}

// ID returns a unique identifier for the function signature.
// A package should only have one function with this particular ID.
func (fs FuncSignature) ID() string {
	return fmt.Sprintf("%s.%s", fs.Receiver.Name, fs.Name)
}

func (fs FuncSignature) String() string {
	var sb strings.Builder

	sb.WriteString("func ")

	// if method, add the receiver to the signature
	if fs.Receiver.IsDefined() {
		sb.WriteString(fmt.Sprintf("(%s) ", fs.Receiver))
	}

	// write function name and parameters
	sb.WriteString(fmt.Sprintf("%s(%s)", fs.Name, fs.Params))

	// add results, if any are specified
	if len(fs.Results) == 1 {
		sb.WriteString(fmt.Sprintf(" %s", fs.Results))
	} else if len(fs.Results) > 1 {
		sb.WriteString(fmt.Sprintf(" (%s)", fs.Results))
	}

	return sb.String()
}

func (fs FuncSignature) compareString() string {
	// TODO: consider underlying package change. For example:
	// before: func(x pkg.Type) depends on module named github.com/a/pkg
	// after: func(x pkg.Type) depends on module named github.com/b/pkg
	return fs.String()
}

// ParseFuncSignature parses a FuncDecl into a FuncSignature.
func ParseFuncSignature(decl *ast.FuncDecl) FuncSignature {
	fs := FuncSignature{
		Name:    decl.Name.Name,
		Params:  extractTypeList(decl.Type.Params),
		Results: extractTypeList(decl.Type.Results),
	}

	recvlist := extractTypeList(decl.Recv)
	if len(recvlist) > 0 {
		fs.Receiver = recvlist[0]
	}

	return fs
}
