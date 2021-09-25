package modface

import "fmt"

// EllipsisTypeIdentifier identifies a pointer type.
type EllipsisTypeIdentifier struct {
	TypeIdentifier
}

func (eti EllipsisTypeIdentifier) String() string {
	return fmt.Sprintf("...%s", eti.TypeIdentifier)
}

// TypeID returns a certain, unmistakable type with full package name qualifiers.
func (eti EllipsisTypeIdentifier) typeID() string {
	return fmt.Sprintf("...%s", eti.TypeIdentifier.typeID())
}
