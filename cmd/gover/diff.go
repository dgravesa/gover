package main

import (
	"fmt"

	"github.com/spf13/cobra"
)

var diffModpaths = []string{"."}

var diffErrors *optset
var diffChanges *optset

var diffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Print module interface diff",
	Long: `
	Print module interface changes since most recent version based on git tags.
	`,
	RunE: func(cmd *cobra.Command, args []string) error {
		changes, err := diffChanges.Value()
		if err != nil {
			return err
		}

		errcond, err := diffErrors.Value()
		if err != nil {
			return err
		}

		if len(args) != 0 {
			diffModpaths = args
		}

		fmt.Println(changes, errcond, diffModpaths)
		return nil
	},
}

func init() {
	diffChanges = makeOptsetFlag(diffCmd.Flags(), "changes", "c", "changes to print",
		"any", "breaking")

	diffErrors = makeOptsetFlag(diffCmd.Flags(), "errors", "e",
		"condition to exit with error status code",
		"none", "any", "breaking")
}
