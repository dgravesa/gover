package modface

// Export represents an export.
type Export interface {
	// String returns a human-readable presentation of the export.
	String() string

	// ID returns a unique identifier for an export.
	// An export's ID must be unique within its package.
	ID() string

	// compareString returns a complete string representation of the export such that any two
	// exports with matching compareStrings may be considered equal, and any two exports with
	// differing compareStrings may be considered not equal.
	compareString() string
}

// ExportsEqual returns true if faces are equal, otherwise false.
func ExportsEqual(a, b Export) bool {
	return a.compareString() == b.compareString()
}
