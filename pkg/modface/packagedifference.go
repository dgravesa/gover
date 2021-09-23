package modface

// PackageDifference returns the interface differences between two versions of a package.
type PackageDifference struct {
	Additions map[string]Export
	Removals  map[string]Export
	Changes   map[string]ExportDifference
}

func newPackageDifference() *PackageDifference {
	pd := new(PackageDifference)
	pd.Additions = make(map[string]Export)
	pd.Removals = make(map[string]Export)
	pd.Changes = make(map[string]ExportDifference)
	return pd
}

// ExportDifference contains the old and new faces.
type ExportDifference struct {
	Old Export
	New Export
}

// Any returns true if there are any differences, otherwise false.
func (pd PackageDifference) Any() bool {
	if len(pd.Additions) > 0 || len(pd.Removals) > 0 || len(pd.Changes) > 0 {
		return true
	}
	return false
}

// Breaking returns true if there are any breaking differences, otherwise false.
// Any interface removals or changes in signature are considered breaking changes.
func (pd PackageDifference) Breaking() bool {
	if len(pd.Removals) > 0 || len(pd.Changes) > 0 {
		return true
	}
	return false
}

// PackageDiff returns an object representing the difference between two package versions.
func PackageDiff(oldpack, newpack PackageInterface) *PackageDifference {
	packdiff := newPackageDifference()

	for id, oldface := range oldpack {
		newface, found := newpack[id]
		if !found {
			// face in old but not in new, so it has been removed
			packdiff.Removals[id] = oldface
		} else if !ExportsEqual(oldface, newface) {
			// face has changed
			packdiff.Changes[id] = ExportDifference{
				Old: oldface,
				New: newface,
			}
		}
	}

	for id, newface := range newpack {
		_, found := oldpack[id]
		if !found {
			// face in new but not in old, so it has been added
			packdiff.Additions[id] = newface
		}
	}

	return packdiff
}
