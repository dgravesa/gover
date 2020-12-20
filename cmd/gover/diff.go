package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/dgravesa/gover/pkg/modface"
	"github.com/dgravesa/minicli"
)

type diffCmd struct {
	modpath  *string // injected by main command
	pchanges *optset
	errcond  *optset
	compare  string
}

func newDiffCmd(modpath *string) minicli.CmdImpl {
	return &diffCmd{modpath: modpath}
}

func (d *diffCmd) SetFlags(flags *flag.FlagSet) {
	d.pchanges = makeOptsetFlag(flags, "changes", "changes to print", "any", "breaking")
	d.errcond = makeOptsetFlag(flags, "error", "condition to exit with error status code",
		"none", "breaking", "any")
	flags.StringVar(&d.compare, "compare", "HEAD", "specify commit or tag to compare against")
}

func (d *diffCmd) Exec(args []string) error {
	// validate changes and error arguments
	pchanges, err := d.pchanges.Value()
	if err != nil {
		return err
	}
	errcond, err := d.errcond.Value()
	if err != nil {
		return err
	}

	modpath := *d.modpath

	var currentModule *modface.Module
	var compareModule *modface.Module
	currentDone := make(chan error)
	recentDone := make(chan error)

	// parse current module interface
	go func() {
		var err error
		currentModule, err = modface.ParseModule(modpath)
		currentDone <- err
	}()

	// copy module and checkout compare version
	go func() {
		// create temporary directory
		tmpdir, err := ioutil.TempDir(os.TempDir(), "gover-*")
		if err != nil {
			recentDone <- err
			return
		}

		// copy module into temporary directory
		cpycmd := exec.Command("cp", "-r", fmt.Sprintf("%s/", modpath), tmpdir)
		_, cperr := cpycmd.Output()

		// NOTE: doing this manually because defer does not seem to like it
		rmdir := func() {
			rmcmd := exec.Command("rm", "-rf", tmpdir)
			_, rmerr := rmcmd.Output()
			// verify rm succeeded
			switch v := rmerr.(type) {
			case nil:
				break
			case *exec.ExitError:
				fmt.Fprintln(os.Stderr, rmcmd.String(), string(v.Stderr))
			default:
				fmt.Fprintln(os.Stderr, "error:", v)
			}
		}

		// verify copy succeeded
		switch v := cperr.(type) {
		case nil:
			break
		case *exec.ExitError:
			rmdir()
			recentDone <- fmt.Errorf("%s: %s", cpycmd, string(v.Stderr))
			return
		default:
			rmdir()
			recentDone <- fmt.Errorf("%s: %v", cpycmd, cperr)
			return
		}

		// checkout recent version
		checkoutcmd := exec.Command("git", "-C", tmpdir, "checkout", "-f", d.compare)
		_, checkouterr := checkoutcmd.Output()
		switch v := checkouterr.(type) {
		case nil:
			break
		case *exec.ExitError:
			rmdir()
			recentDone <- fmt.Errorf("%s: %s", checkoutcmd, string(v.Stderr))
			return
		default:
			rmdir()
			recentDone <- fmt.Errorf("%s: %v", checkoutcmd, err)
			return
		}

		// parse version of module
		compareModule, err = modface.ParseModule(tmpdir)
		rmdir()
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
	moduleDifference := modface.Diff(compareModule, currentModule)

	// print differences to stdout as specified by change level
	printDiff(currentModule.Path, moduleDifference, pchanges)

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

func printDiff(modname string, moduleDifference *modface.ModuleDifference, level string) {
	type difference interface {
		Any() bool
		Breaking() bool
	}

	meetsLevel := func(d difference) bool {
		if level == "any" && d.Any() {
			return true
		} else if level == "breaking" && d.Breaking() {
			return true
		}
		return false
	}

	if !meetsLevel(moduleDifference) {
		return
	}

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
		if !meetsLevel(pkgchanges) {
			continue
		}

		fmt.Println("---", "package", pkgname)

		// print package removals
		for _, face := range pkgchanges.Removals {
			fmt.Println("<  ", face)
		}
		if level == "any" {
			// print package additions
			for _, face := range pkgchanges.Additions {
				fmt.Println(">  ", face)
			}
		}
		// print package changes
		for _, facediff := range pkgchanges.Changes {
			fmt.Println("<  ", facediff.Old)
			fmt.Println(">  ", facediff.New)
		}
	}
}
