package modface

import (
	"io/ioutil"
	"path/filepath"

	"github.com/dgravesa/gover/pkg/modparse"
	"golang.org/x/mod/modfile"
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
		parseDir(module.Packages, dir, module.Path)
	}

	return module, nil
}
