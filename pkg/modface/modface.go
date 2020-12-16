package modface

// ModFace represents all exports of a module.
type ModFace map[string]PackFace

// Face represents an export.
type Face interface {
	String() string
	ID() string
}

// ParseMod parses a module and returns all of its export signatures.
func ParseMod(moddir string) (ModFace, error) {
	mf := make(ModFace)

	// TODO: do this for all dirs in the module
	parseDir(moddir, mf)

	return mf, nil
}
