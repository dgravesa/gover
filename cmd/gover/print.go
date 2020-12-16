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
			err := printModpath(modpath)
			if err != nil {
				return err
			}
		}

		return nil
	},
}

func printModpath(moddir string) error {
	mf, err := modface.ParseMod(moddir)
	if err != nil {
		return err
	}

	for pkgname, pkgface := range mf {
		fmt.Println("--------------------")
		fmt.Println("package", pkgname)
		fmt.Println("--------------------")

		for _, face := range pkgface {
			fmt.Println(face)
		}

		fmt.Println()
	}

	return nil
}
