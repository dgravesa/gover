package main

import (
	"flag"
	"fmt"

	"github.com/dgravesa/minicli"

	"github.com/dgravesa/gover/pkg/modface"
)

type printCmd struct {
	moddir *string
}

func newPrintCmd(moddir *string) minicli.CmdImpl {
	return &printCmd{moddir}
}

func (p *printCmd) SetFlags(_ *flag.FlagSet) {
	// no flags to set
}

func (p *printCmd) Exec(args []string) error {
	module, err := modface.ParseModule(*p.moddir)
	if err != nil {
		return err
	}

	fmt.Println("module", module.Path)

	for pkgname, pkgface := range module.Packages {
		fmt.Println("- package", pkgname)

		for _, face := range pkgface {
			fmt.Println("  -", face)
		}
	}

	return nil
}
