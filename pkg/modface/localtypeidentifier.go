package modface

import (
	"fmt"
)

// LocalTypeIdentifier corresponds to either a type defined within the package or a built-in type.
// TODO: should locally-defined be grouped with imported types instead?
type LocalTypeIdentifier struct {
	TypeName  string
	IsPointer bool
}

func (lti LocalTypeIdentifier) String() string {
	if lti.IsPointer {
		return fmt.Sprintf("*%s", lti.TypeName)
	}
	return lti.TypeName
}

// TypeID returns a unique ID for the type across any go packages.
func (lti LocalTypeIdentifier) TypeID() string {
	return lti.String()
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
