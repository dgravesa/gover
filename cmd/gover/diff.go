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

		for _, modpath := range diffModpaths {
			err := printModuleInterfaceDiff(modpath, changes, errcond)
			if err != nil {
				return err
			}
		}

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

func printModuleInterfaceDiff(modpath string, pchanges, errcond string) error {
	var currentModule *modface.Module
	var recentVersionModule *modface.Module
	var tmpdir string
	var mostRecentVersion string
	currentDone := make(chan error)
	copyDone := make(chan error)
	versionsDone := make(chan error)
	recentDone := make(chan error)

	// parse current module interface
	go func() {
		var err error
		currentModule, err = modface.ParseModule(modpath)
		currentDone <- err
	}()

	// make copy of module
	go func() {
		var err error
		tmpdir, err = ioutil.TempDir(os.TempDir(), "gover-*")
		if err != nil {
			copyDone <- err
			return
		}

		// copy module into temporary directory
		cpycmd := exec.Command("cp", "-r", fmt.Sprintf("%s/", modpath), tmpdir)
		_, cperr := cpycmd.Output()
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
		copyDone <- err
	}()

	// get versions
	go func() {
		versions, err := modface.Versions(modpath)
		if err != nil {
			versionsDone <- err
			return
		} else if len(versions) == 0 {
			versionsDone <- fmt.Errorf("%s : no versions detected", modpath)
			return
		}

		// get most recent version
		mostRecentVersion = versions[0]
		for _, version := range versions {
			mostRecentVersion = semver.Max(version, mostRecentVersion)
		}

		versionsDone <- nil
	}()

	// parse recent version module interface
	go func() {
		copyErr := <-copyDone
		if copyErr != nil {
			recentDone <- copyErr
			return
		}

		rmcmd := exec.Command("rm", "-rf", tmpdir)
		defer rmcmd.Run()

		versionsErr := <-versionsDone
		if versionsErr != nil {
			recentDone <- versionsErr
			return
		}

		// checkout recent version
		checkoutcmd := exec.Command("git", "-C", tmpdir, "checkout", mostRecentVersion)
		_, err := checkoutcmd.Output()
		switch v := err.(type) {
		case nil:
			break
		case *exec.ExitError:
			recentDone <- fmt.Errorf("%s: %s", checkoutcmd, string(v.Stderr))
			return
		default:
			recentDone <- fmt.Errorf("%s: %v", checkoutcmd, err)
			return
		}

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

	if moduleDifference.Any() {
		// print module name
		fmt.Println("module", currentModule.Path)
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

	return nil
}
