package modface

// ModuleDifference represents the interface difference between two versions of a module.
type ModuleDifference struct {
	PackageRemovals  map[string]PackageInterface
	PackageAdditions map[string]PackageInterface
	PackageChanges   map[string]*PackageDifference
}

func newModuleDifference() *ModuleDifference {
	md := new(ModuleDifference)
	md.PackageRemovals = make(map[string]PackageInterface)
	md.PackageAdditions = make(map[string]PackageInterface)
	md.PackageChanges = make(map[string]*PackageDifference)
	return md
}

// Any returns true if there are any differences, otherwise false.
func (md ModuleDifference) Any() bool {
	if len(md.PackageAdditions) > 0 || len(md.PackageRemovals) > 0 || len(md.PackageChanges) > 0 {
		return true
	}
	return false
}

// Breaking returns true if there are any breaking differences, otherwise false.
// Any package removals or packages with breaking changes are considered breaking changes
// for the module.
func (md ModuleDifference) Breaking() bool {
	if len(md.PackageRemovals) > 0 {
		return true
	}
	for _, packdiff := range md.PackageChanges {
		if packdiff.Breaking() {
			return true
		}
	}
	return false
}

// Diff computes the interface difference between two versions of a module.
func Diff(oldmod, newmod *Module) *ModuleDifference {
	moddiff := newModuleDifference()

	for pkgname, oldpack := range oldmod.Packages {
		newpack, found := newmod.Packages[pkgname]
		if !found {
			// package in old but not in new, so it has been removed
			moddiff.PackageRemovals[pkgname] = oldpack
		} else {
			packdiff := PackageDiff(oldpack, newpack)
			if packdiff.Any() {
				moddiff.PackageChanges[pkgname] = packdiff
			}
		}
	}

	for pkgname, newface := range newmod.Packages {
		_, found := oldmod.Packages[pkgname]
		if !found {
			// package in new but not in old, so it has been added
			moddiff.PackageAdditions[pkgname] = newface
		}
	}

	return moddiff
}
