package modface

// TypeIdentifier is an interface to return a complete type ID, including full package name.
type TypeIdentifier interface {
	TypeID() string
}
