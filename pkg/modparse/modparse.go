package modparse

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ModDirs returns relative paths to all directories which may be considered a part of the module.
func ModDirs(path string) ([]string, error) {
	return modDirs(path, true)
}

func modDirs(path string, thismod bool) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	dir, err := f.Readdir(0)
	if err != nil {
		return nil, err
	}

	moddirs := []string{path}
	ismod := false

	isValidSubdir := func(fi os.FileInfo) bool {
		if !fi.IsDir() {
			return false
		} else if strings.HasPrefix(fi.Name(), ".") {
			return false
		} else if fi.Name() == "vendor" {
			return false
		}
		return true
	}

	for _, fi := range dir {
		if fi.Name() == "go.mod" && !fi.IsDir() {
			// go.mod file detected
			if thismod {
				ismod = true
			} else {
				// subdirectory contains a different module
				return []string{}, nil
			}
		} else if isValidSubdir(fi) {
			relpath := filepath.Join(path, fi.Name())
			// module subdirectory
			subdirs, err := modDirs(relpath, false)
			if err != nil {
				return nil, err
			}
			moddirs = append(moddirs, subdirs...)
		}
	}

	if thismod && !ismod {
		return nil, fmt.Errorf("%s does not define a go module", path)
	}

	return moddirs, nil
}
