package modface

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// FuncSig defines a function signature.
// The ID is a short signature that should be uniquely identifying.
// The Signature is a complete representation of the function's interface
// and should be directly comparable between different commits to ensure
// that backwards compatibility is maintained.
type FuncSig struct {
	id        string
	signature string
}

// ID returns a unique identifier for the function signature.
// A package should only have one function with this particular ID.
func (fs FuncSig) ID() string {
	return fs.id
}

func (fs FuncSig) String() string {
	return fs.signature
}

// ParseFuncSig parses a FuncDecl into a FuncSig.
func ParseFuncSig(decl *ast.FuncDecl) FuncSig {
	funcname := decl.Name.Name
	recvtypes := flTypes(decl.Recv)
	paramtypes := flTypes(decl.Type.Params)
	resulttypes := flTypes(decl.Type.Results)

	var recvstr string
	var paramsstr string
	var resultsstr string

	// build recv string
	if len(recvtypes) > 0 {
		recvstr = fmt.Sprintf("(%s) ", strings.Join(recvtypes, ", "))
	} else {
		recvstr = ""
	}

	// build params string
	paramsstr = fmt.Sprintf("(%s)", strings.Join(paramtypes, ", "))

	// build results string
	if len(resulttypes) == 0 {
		resultsstr = ""
	} else if len(resulttypes) == 1 {
		resultsstr = fmt.Sprintf(" %s", resulttypes[0])
	} else {
		resultsstr = fmt.Sprintf(" (%s)", strings.Join(resulttypes, ", "))
	}

	return FuncSig{
		id:        fmt.Sprintf("%s%s", recvstr, funcname),
		signature: fmt.Sprintf("func %s%s%s%s", recvstr, funcname, paramsstr, resultsstr),
	}
}

func flTypes(fl *ast.FieldList) []string {
	fltypes := []string{}
	if fl != nil {
		for _, f := range fl.List {
			str := typeStr(f.Type)
			fltypes = appendStrN(fltypes, str, maxInt(1, len(f.Names)))
		}
	}
	return fltypes
}

func typeStr(t ast.Expr) string {
	var str string
	switch v := t.(type) {
	case *ast.FuncType:
		str = funcTypeStr(v)
	default:
		str = types.ExprString(v)
	}
	return str
}

func funcTypeStr(f *ast.FuncType) string {
	if f == nil {
		return ""
	}

	params := []string{}
	results := []string{}

	// get param type strings
	if f.Params != nil {
		for _, p := range f.Params.List {
			typestr := typeStr(p.Type)
			params = appendStrN(params, typestr, maxInt(1, len(p.Names)))
		}
	}
	// get result type strings
	if f.Results != nil {
		for _, r := range f.Results.List {
			typestr := typeStr(r.Type)
			results = appendStrN(results, typestr, maxInt(1, len(r.Names)))
		}
	}

	// construct signature
	if len(results) == 0 {
		return fmt.Sprintf("func(%s)", strings.Join(params, ", "))
	} else if len(results) == 1 {
		return fmt.Sprintf("func(%s)%s", strings.Join(params, ", "), results[0])
	}
	return fmt.Sprintf("func(%s)(%s)", strings.Join(params, ", "), strings.Join(results, ", "))
}

func appendStrN(slice []string, str string, N int) []string {
	strs := []string{}
	for i := 0; i < N; i++ {
		strs = append(strs, str)
	}
	return append(slice, strs...)
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
