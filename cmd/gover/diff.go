package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/dgravesa/gover/pkg/modface"
	"github.com/spf13/cobra"
	"golang.org/x/mod/semver"
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
	Run: func(cmd *cobra.Command, args []string) {
		changes, err := diffChanges.Value()
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		errcond, err := diffErrors.Value()
		if err != nil {
			fmt.Println("error:", err)
			os.Exit(1)
		}

		if len(args) != 0 {
			diffModpaths = args
		}

		errorOccurred := false
		for _, modpath := range diffModpaths {
			err := printModuleInterfaceDiff(modpath, changes, errcond)
			if err != nil {
				fmt.Println("error:", err)
				errorOccurred = true
			}
		}

		if errorOccurred {
			os.Exit(2)
		}
	},
}

func init() {
	diffChanges = makeOptsetFlag(diffCmd.Flags(), "changes", "c", "changes to print",
		"any", "breaking")

	diffErrors = makeOptsetFlag(diffCmd.Flags(), "errors", "e",
		"condition to exit with error status code",
		"none", "any", "breaking")
}

func printModuleInterfaceDiff(modpath string, pchanges, errcond string) error {
	var currentModule *modface.Module
	var recentVersionModule *modface.Module
	currentDone := make(chan error)
	recentDone := make(chan error)

	// parse current module interface
	go func() {
		var err error
		currentModule, err = modface.ParseModule(modpath)
		currentDone <- err
	}()

	// copy module and checkout recent version
	go func() {
		versions, err := modface.Versions(modpath)
		if err != nil {
			recentDone <- err
			return
		} else if len(versions) == 0 {
			recentDone <- fmt.Errorf("%s : no versions detected", modpath)
			return
		}

		// get most recent version
		mostRecentVersion := versions[0]
		for _, version := range versions {
			mostRecentVersion = semver.Max(version, mostRecentVersion)
		}

		// create temporary directory
		tmpdir, err := ioutil.TempDir(os.TempDir(), "gover-*")
		if err != nil {
			recentDone <- err
			return
		}

		// copy module into temporary directory
		cpycmd := exec.Command("cp", "-r", fmt.Sprintf("%s/", modpath), tmpdir)
		_, cperr := cpycmd.Output()
		rmcmd := exec.Command("rm", "-rf", tmpdir)
		defer rmcmd.Run()

		// verify copy succeeded
		switch v := cperr.(type) {
		case nil:
			break
		case *exec.ExitError:
			recentDone <- fmt.Errorf("%s: %s", cpycmd, string(v.Stderr))
			return
		default:
			recentDone <- fmt.Errorf("%s: %v", cpycmd, cperr)
			return
		}

		// checkout recent version
		checkoutcmd := exec.Command("git", "-C", tmpdir, "checkout", mostRecentVersion)
		_, checkouterr := checkoutcmd.Output()
		switch v := checkouterr.(type) {
		case nil:
			break
		case *exec.ExitError:
			recentDone <- fmt.Errorf("%s: %s", checkoutcmd, string(v.Stderr))
			return
		default:
			recentDone <- fmt.Errorf("%s: %v", checkoutcmd, err)
			return
		}

		// parse version of module
		recentVersionModule, err = modface.ParseModule(tmpdir)
		recentDone <- err
	}()

	// wait for results
	currentErr := <-currentDone
	recentErr := <-recentDone

	if currentErr != nil {
		return currentErr
	} else if recentErr != nil {
		return recentErr
	}

	// compute difference between module interfaces
	moduleDifference := modface.Diff(recentVersionModule, currentModule)

	switch pchanges {
	case "breaking":
		if moduleDifference.Breaking() {
			printBreaking(currentModule.Path, moduleDifference)
		}
	default:
		if moduleDifference.Any() {
			printAny(currentModule.Path, moduleDifference)
		}
	}

	var resultStatus error
	switch errcond {
	case "breaking":
		if moduleDifference.Breaking() {
			resultStatus = fmt.Errorf("breaking changes detected")
		}
	case "any":
		if moduleDifference.Any() {
			resultStatus = fmt.Errorf("changes detected")
		}
	default:
		resultStatus = nil
	}

	return resultStatus
}

func printAny(modname string, moduleDifference *modface.ModuleDifference) {
	// print module name
	fmt.Println("module", modname)
	// print removals
	for pkgname := range moduleDifference.PackageRemovals {
		fmt.Println("< package", pkgname)
	}
	// print additions
	for pkgname := range moduleDifference.PackageAdditions {
		fmt.Println("> package", pkgname)
	}
	// print changes per package
	for pkgname, pkgchanges := range moduleDifference.PackageChanges {
		fmt.Println("- package", pkgname)

		// print package removals
		for _, face := range pkgchanges.Removals {
			fmt.Println("  <", face)
		}
		// print package additions
		for _, face := range pkgchanges.Additions {
			fmt.Println("  >", face)
		}
		// print package changes
		for _, facediff := range pkgchanges.Changes {
			fmt.Println("  -", facediff.Old, "->", facediff.New)
		}
	}
}

func printBreaking(modname string, moduleDifference *modface.ModuleDifference) {
	// print module name
	fmt.Println("module", modname)
	// print removals
	for pkgname := range moduleDifference.PackageRemovals {
		fmt.Println("< package", pkgname)
	}
	// print changes per package
	for pkgname, pkgchanges := range moduleDifference.PackageChanges {
		if pkgchanges.Breaking() {
			fmt.Println("- package", pkgname)

			// print package removals
			for _, face := range pkgchanges.Removals {
				fmt.Println("  <", face)
			}
			// print package changes
			for _, facediff := range pkgchanges.Changes {
				fmt.Println("  -", facediff.Old, "->", facediff.New)
			}
		}
	}
}
