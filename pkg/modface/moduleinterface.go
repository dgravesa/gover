package modface

import (
	"io/ioutil"
	"path/filepath"
	"strings"

	"os/exec"

	"github.com/dgravesa/gover/pkg/modparse"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

// ModuleInterface represents all exports of a module.
type ModuleInterface map[string]PackageInterface

// Module represents a module.
// The Package member contains all exports of the module.
type Module struct {
	Path     string
	Packages map[string]PackageInterface
}

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

// ParseModule parses a module and returns all of its export signatures.
func ParseModule(moddir string) (*Module, error) {
	dirs, err := modparse.ModuleDirs(moddir)
	if err != nil {
		return nil, err
	}

	mfpath := filepath.Join(moddir, "go.mod")
	mfile, err := ioutil.ReadFile(mfpath)
	if err != nil {
		return nil, err
	}

	module := new(Module)
	module.Path = modfile.ModulePath(mfile)
	module.Packages = make(ModuleInterface)

	for _, dir := range dirs {
		parseDir(module.Packages, moddir, dir, module.Path)
	}

	return module, nil
}

// Versions returns all the versions for a module pointed to by moddir.
func Versions(moddir string) ([]string, error) {
	tagcmd := exec.Command("git", "-C", moddir, "tag")
	tagout, err := tagcmd.Output()
	if err != nil {
		return nil, err
	}

	tags := strings.Split(string(tagout), "\n")

	versions := []string{}
	for _, tag := range tags {
		if semver.IsValid(tag) {
			versions = append(versions, tag)
		}
	}

	return versions, nil
}
