package modface

import (
	"fmt"
)

// StructSignature defines a struct type signature.
type StructSignature struct {
	Name        string
	PackageName string

	// TODO: maybe this?
	Fields []FieldSignature
}

// ID returns a unique identifier for the type within the package.
func (ts StructSignature) ID() string {
	return ts.Name
}

// TypeID returns a unique identifier for the type among all packages.
func (ts StructSignature) TypeID() string {
	return fmt.Sprintf("%s.%s", ts.PackageName, ts.Name)
}

func (ts StructSignature) String() string {
	// TODO: implement
	return ts.Name
}

func (ts StructSignature) compareString() string {
	// TODO: implement
	return ts.Name
}
