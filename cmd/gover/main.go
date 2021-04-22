package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/dgravesa/minicli"
)

func main() {
	var modpath string

	cli := minicli.New()

	cli.Flags("", "", func(flags *flag.FlagSet) {
		flags.StringVar(&modpath, "C", ".", "path to module")
	})

	cli.Func("print", "print module interface", makePrintFunc(&modpath))

	cli.Cmd("diff", "compare module interface changes to previous version",
		newDiffCmd(&modpath))

	// TODO: cowardly removing for now, needs more work to be safer
	// cli.Cmd("tag", "tag with a suggested version", newTagCmd(&modpath))

	cli.Func("suggest", "suggest a new semantic version", makeSuggestFunc(&modpath))

	if err := cli.Exec(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
