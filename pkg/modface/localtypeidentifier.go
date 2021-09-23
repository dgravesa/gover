package modface

// LocalTypeIdentifier corresponds to either a type defined within the package or a built-in type.
type LocalTypeIdentifier struct {
	TypeName string
}

func (lti LocalTypeIdentifier) String() string {
	return lti.TypeName
}

// TypeID returns a unique ID for the type across any go packages.
func (lti LocalTypeIdentifier) TypeID() string {
	return lti.String()
}
