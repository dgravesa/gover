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

// Face represents an export.
type Face interface {
	String() string
	ID() string
}

// FacesEqual returns true if faces are equal, otherwise false.
func FacesEqual(a, b Face) bool {
	return a.String() == b.String()
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
