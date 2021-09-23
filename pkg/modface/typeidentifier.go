package modface

// TypeIdentifier is an interface to return both a human-readable representation of a type
// and a certain version of the type with full package name qualifiers.
type TypeIdentifier interface {
	// String returns a human-readable representation of a type.
	String() string

	// TypeID returns a certain, unmistakable type with full package name qualifiers.
	TypeID() string
}
