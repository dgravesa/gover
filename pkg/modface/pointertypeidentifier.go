package modface

import "fmt"

// PointerTypeIdentifier identifies a pointer type.
type PointerTypeIdentifier struct {
	TypeIdentifier
}

func (pti PointerTypeIdentifier) String() string {
	return fmt.Sprintf("*%s", pti.TypeIdentifier)
}

// TypeID returns a certain, unmistakable type with full package name qualifiers.
func (pti PointerTypeIdentifier) typeID() string {
	return fmt.Sprintf("*%s", pti.TypeIdentifier.typeID())
}
