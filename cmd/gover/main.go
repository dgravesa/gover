package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dgravesa/minicli"
)

func main() {
	var modpath string

	minicli.Flags("", "", func(flags *flag.FlagSet) {
		flags.StringVar(&modpath, "C", ".", "path to module")
	})

	minicli.Func("print", "print module interface", makePrintFunc(&modpath))

	minicli.Cmd("diff", "compare module interface changes to previous version",
		newDiffCmd(&modpath))

	if err := minicli.Exec(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
