package modface

import "fmt"

// SliceTypeIdentifier identifies a slice type.
type SliceTypeIdentifier struct {
	ElementType TypeIdentifier
}

func (sti SliceTypeIdentifier) String() string {
	return fmt.Sprintf("[]%s", sti.ElementType)
}

// TypeID returns a certain, unmistakable type with full package name qualifiers.
func (sti SliceTypeIdentifier) typeID() string {
	return fmt.Sprintf("[]%s", sti.ElementType.typeID())
}
