package main

import (
	"fmt"

	"github.com/dgravesa/gover/pkg/modface"
	"github.com/spf13/cobra"
)

var printModpaths = []string{"."}

var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Print module interface",
	Long: `
	Print all exports of a module
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			printModpaths = args
		}

		for _, modpath := range printModpaths {
			err := printModuleInterface(modpath)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func printModuleInterface(moddir string) error {
	module, err := modface.ParseModule(moddir)
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
