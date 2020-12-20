package main

import (
	"fmt"

	"github.com/dgravesa/gover/pkg/modface"
)

func makePrintFunc(moddir *string) func(args []string) error {
	return func(args []string) error {
		module, err := modface.ParseModule(*moddir)
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
}
