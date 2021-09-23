package modface

import (
	"fmt"
	"go/ast"
	"path/filepath"
	"strings"
)

// BaseTypeIdentifier is the type signature corresponding to types whose underlying details are not
// relevant to the package interface aside from unique identification of the type across any
// go packages.
type BaseTypeIdentifier struct {
	PackageName string
	TypeName    string
	IsPointer   bool
}

func (bs BaseTypeIdentifier) String() string {
	pkgShortName := filepath.Base(bs.PackageName)
	if pkgShortName != "." {
		return fmt.Sprintf("%s.%s", pkgShortName, bs.TypeName)
	}
	return bs.TypeName
}

// TypeID returns a unique ID for the type across any go packages.
func (bs BaseTypeIdentifier) TypeID() string {
	return fmt.Sprintf("%s.%s", bs.PackageName, bs.TypeName)
}

func parseBaseTypeIdentifier(v *ast.Ident, cache *importCache) (BaseTypeIdentifier, error) {
	name := strings.TrimPrefix(v.Name, "*")
	isPointer := strings.HasPrefix(v.Name, "*")

	nameSplit := strings.Split(name, ".")
	fmt.Println("NAMESPLIT: ", nameSplit) // TODO: remove
	if len(nameSplit) == 1 {
		// return as local type or primitive
		return BaseTypeIdentifier{
			TypeName:  name,
			IsPointer: isPointer,
		}, nil
	} else if len(nameSplit) != 2 {
		// TODO: better error typing
		return BaseTypeIdentifier{}, fmt.Errorf("invalid name: %s", v.Name)
	}

	importShortName := nameSplit[0]
	typeName := nameSplit[1]

	// find full import package name from short identifier
	fullImportName, found := cache.FindShortName(importShortName)
	if !found {
		// TODO: better error typing
		return BaseTypeIdentifier{}, fmt.Errorf("import not found: %s", importShortName)
	}

	return BaseTypeIdentifier{
		PackageName: fullImportName,
		TypeName:    typeName,
		IsPointer:   isPointer,
	}, nil
}
