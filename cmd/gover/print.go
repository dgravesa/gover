package main

import (
	"fmt"
	"strings"

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

	printHeader := func(header, linechar string) {
		linestr := strings.Repeat(linechar, len(header))
		fmt.Println(linestr)
		fmt.Println(header)
		fmt.Println(linestr)
	}

	printHeader("module "+module.Path, "=")
	fmt.Println()

	for pkgname, pkgface := range module.Packages {
		printHeader("package "+pkgname, "-")

		for _, face := range pkgface {
			fmt.Println(face)
		}

		fmt.Println()
	}

	return nil
}
