package modface

import (
	"fmt"
	"path/filepath"
)

// ImportedTypeIdentifier corresponds to a type imported from another package.
type ImportedTypeIdentifier struct {
	PackageName string
	TypeName    string
	IsPointer   bool
}

func (iti ImportedTypeIdentifier) String() string {
	pkgShortName := filepath.Base(iti.PackageName)
	if iti.IsPointer {
		return fmt.Sprintf("*%s.%s", pkgShortName, iti.TypeName)
	}
	return fmt.Sprintf("%s.%s", pkgShortName, iti.TypeName)
}

// TypeID returns a unique ID for the type across any go packages.
func (iti ImportedTypeIdentifier) TypeID() string {
	if iti.IsPointer {
		return fmt.Sprintf("*%s.%s", iti.PackageName, iti.TypeName)
	}
	return fmt.Sprintf("%s.%s", iti.PackageName, iti.TypeName)
}

// func parseBaseTypeIdentifier(v *ast.Ident, cache *importCache) (BaseTypeIdentifier, error) {
// 	name := strings.TrimPrefix(v.Name, "*")
// 	isPointer := strings.HasPrefix(v.Name, "*")

// 	nameSplit := strings.Split(name, ".")
// 	fmt.Println("NAMESPLIT: ", nameSplit) // TODO: remove
// 	if len(nameSplit) == 1 {
// 		// return as local type or primitive
// 		return BaseTypeIdentifier{
// 			TypeName:  name,
// 			IsPointer: isPointer,
// 		}, nil
// 	} else if len(nameSplit) != 2 {
// 		// TODO: better error typing
// 		return BaseTypeIdentifier{}, fmt.Errorf("invalid name: %s", v.Name)
// 	}

// 	importShortName := nameSplit[0]
// 	typeName := nameSplit[1]

// 	// find full import package name from short identifier
// 	fullImportName, found := cache.FindShortName(importShortName)
// 	if !found {
// 		// TODO: better error typing
// 		return BaseTypeIdentifier{}, fmt.Errorf("import not found: %s", importShortName)
// 	}

// 	return BaseTypeIdentifier{
// 		PackageName: fullImportName,
// 		TypeName:    typeName,
// 		IsPointer:   isPointer,
// 	}, nil
// }
