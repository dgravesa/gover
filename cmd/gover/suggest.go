package main

import (
	"fmt"
)

func makeSuggestFunc(modpath *string) func(args []string) error {
	return func(args []string) error {
		suggestedVersion := "v0.1.0"
		modpath := *modpath

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
				suggestedVersion = suggestVersion(latestVersion, "breaking")
			} else if moduleDifference.Any() {
				suggestedVersion = suggestVersion(latestVersion, "feature")
			} else {
				suggestedVersion = suggestVersion(latestVersion, "bugfix")
			}
		}

		fmt.Println(suggestedVersion)

		return err
	}
}
