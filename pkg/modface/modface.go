package modface

import (
	"io/ioutil"
	"path/filepath"

	"github.com/dgravesa/gover/pkg/modparse"
	"golang.org/x/mod/modfile"
)

// ModFace represents all exports of a module.
type ModFace map[string]PackFace

// Face represents an export.
type Face interface {
	String() string
	ID() string
}

// ParseMod parses a module and returns all of its export signatures.
func ParseMod(moddir string) (ModFace, error) {
	dirs, err := modparse.ModDirs(moddir)
	if err != nil {
		return nil, err
	}

	mfpath := filepath.Join(moddir, "go.mod")
	mf, err := ioutil.ReadFile(mfpath)
	if err != nil {
		return nil, nil
	}
	modname := modfile.ModulePath(mf)

	mface := make(ModFace)
	for _, dir := range dirs {
		parseDir(mface, dir, modname)
	}

	return mface, nil
}
