package modface

import (
	"io/ioutil"
	"path/filepath"

	"github.com/dgravesa/gover/pkg/modparse"
	"golang.org/x/mod/modfile"
)

// ModuleInterface represents all exports of a module.
type ModuleInterface map[string]PackageInterface

// Face represents an export.
type Face interface {
	String() string
	ID() string
}

// ParseModule parses a module and returns all of its export signatures.
func ParseModule(moddir string) (ModuleInterface, error) {
	dirs, err := modparse.ModuleDirs(moddir)
	if err != nil {
		return nil, err
	}

	mfpath := filepath.Join(moddir, "go.mod")
	mf, err := ioutil.ReadFile(mfpath)
	if err != nil {
		return nil, nil
	}
	modname := modfile.ModulePath(mf)

	mface := make(ModuleInterface)
	for _, dir := range dirs {
		parseDir(mface, dir, modname)
	}

	return mface, nil
}
