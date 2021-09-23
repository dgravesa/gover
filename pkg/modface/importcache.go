package modface

import (
	"go/ast"
	"path/filepath"
	"strings"
)

type importCache struct {
	importsByShortName map[string]string
}

func newImportCache() *importCache {
	return &importCache{
		importsByShortName: map[string]string{},
	}
}

func (ic *importCache) Insert(spec *ast.ImportSpec) {
	fullPath := strings.Trim(spec.Path.Value, "\"")
	var shortName string
	if spec.Name != nil {
		shortName = spec.Name.Name
	} else {
		shortName = filepath.Base(fullPath)
	}
	ic.importsByShortName[shortName] = fullPath
}

func (ic importCache) FindShortName(shortName string) (string, bool) {
	longName, found := ic.importsByShortName[shortName]
	return longName, found
}
