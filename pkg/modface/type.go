package modface

import (
	"fmt"
	"go/ast"
	"go/types"
	"strings"
)

// Type defines a type that can be used for params, results, and receivers.
type Type struct {
	Name      string
	IsPointer bool
}

func (t Type) String() string {
	if t.IsPointer {
		return fmt.Sprintf("*%s", t.Name)
	}
	return t.Name
}

// IsExported returns true if the type is exported, otherwise false.
func (t Type) IsExported() bool {
	return ast.IsExported(t.Name)
}

// IsDefined returns true if the type is not defined, otherwise false.
func (t Type) IsDefined() bool {
	return t.Name != ""
}

// TypeList is a list of Type instances for function params and results.
type TypeList []Type

func (tl TypeList) String() string {
	typestrings := []string{}
	for _, t := range tl {
		typestrings = append(typestrings, t.String())
	}
	return strings.Join(typestrings, ", ")
}

func extractTypeList(fl *ast.FieldList) TypeList {
	types := []Type{}

	fltypes := flTypes(fl)
	for _, fltype := range fltypes {
		var t Type
		if fltype[0] == '*' {
			t.Name = fltype[1:]
			t.IsPointer = true
		} else {
			t.Name = fltype
			t.IsPointer = false
		}
		types = append(types, t)
	}

	return types
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
