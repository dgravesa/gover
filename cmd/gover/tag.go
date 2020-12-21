package main

import (
	"flag"
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"github.com/dgravesa/minicli"
	"golang.org/x/mod/semver"
)

type tagCmd struct {
	modpath    *string // injected by main command
	pushRemote string
	dryRun     bool
	message    string
}

func newTagCmd(modpath *string) minicli.CmdImpl {
	return &tagCmd{
		modpath: modpath,
	}
}

func (tc *tagCmd) SetFlags(flags *flag.FlagSet) {
	flags.StringVar(&tc.pushRemote, "push", "", "specify a remote to push tag")
	flags.BoolVar(&tc.dryRun, "n", false, "print git commands but do not run them")
	flags.StringVar(&tc.message, "m", "", "specify a message for the tag")
}

func (tc *tagCmd) Exec(args []string) error {
	newVersion := "v0.1.0"
	modpath := *tc.modpath

	versions, err := listVersions(modpath)
	if err != nil {
		return err
	}

	if len(versions) > 0 {
		// determine an appropriate next version based on module differences
		latestVersion := versions[len(versions)-1]
		moduleDifference, err := diff(modpath, latestVersion)
		if err != nil {
			return err
		}

		if moduleDifference.Breaking() {
			newVersion = suggestVersion(latestVersion, "breaking")
		} else if moduleDifference.Any() {
			newVersion = suggestVersion(latestVersion, "feature")
		} else {
			newVersion = suggestVersion(latestVersion, "bugfix")
		}
	}

	runcmd := func(name string, args ...string) error {
		cmd := exec.Command(name, args...)
		if tc.dryRun {
			// print command and exit
			fmt.Println(cmd)
			return nil
		}
		_, err := cmd.Output()
		switch v := err.(type) {
		case *exec.ExitError:
			return fmt.Errorf("%s", string(v.Stderr))
		default:
			return v
		}
	}

	// create version tag
	if tc.message == "" {
		tc.message = fmt.Sprintf("version %s", strings.TrimPrefix(newVersion, "v"))
	}
	err = runcmd("git", "-C", modpath, "tag", "-a", newVersion, "-m", tc.message)
	if err != nil {
		return err
	}

	if tc.pushRemote != "" {
		// push version tag to remote
		err = runcmd("git", "-C", modpath, "push", tc.pushRemote, newVersion)
	}

	return err
}

func listVersions(dir string) ([]string, error) {
	tagcmd := exec.Command("git", "-C", dir, "tag")
	tagout, err := tagcmd.Output()
	if err != nil {
		switch v := err.(type) {
		case *exec.ExitError:
			return nil, fmt.Errorf("error: %s", string(v.Stderr))
		default:
			return nil, v
		}
	}

	tags := strings.Split(string(tagout), "\n")

	versions := []string{}
	for _, tag := range tags {
		if semver.IsValid(tag) {
			versions = append(versions, tag)
		}
	}

	sort.SliceStable(versions, func(i, j int) bool {
		return semver.Compare(versions[i], versions[j]) < 0
	})

	return versions, nil
}

func suggestVersion(current string, changeType string) string {
	vfmt := "v%d.%d.%d"
	var major, minor, patch int
	fmt.Sscanf(current, vfmt, &major, &minor, &patch)

	switch changeType {
	case "breaking":
		if major == 0 {
			return fmt.Sprintf(vfmt, 0, minor+1, 0)
		}
		return fmt.Sprintf(vfmt, major+1, 0, 0)
	case "feature":
		if major == 0 {
			return fmt.Sprintf(vfmt, 0, minor, patch+1)
		}
		return fmt.Sprintf(vfmt, major, minor+1, 0)
	case "bugfix":
		if major == 0 {
			return fmt.Sprintf(vfmt, 0, minor, patch+1)
		}
		return fmt.Sprintf(vfmt, major, minor, patch+1)
	}

	panic(fmt.Sprintf("unexpected change type specified: %s", changeType))
}
