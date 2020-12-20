package main

import (
	"flag"
	"fmt"
	"strings"
)

type optset struct {
	valid map[string]struct{}
	name  string
	value string
}

func makeOptsetFlag(flags *flag.FlagSet, name, usage, dflt string, valid ...string) *optset {
	os := &optset{
		valid: map[string]struct{}{
			dflt: {},
		},
		name: name,
	}
	for _, v := range valid {
		os.valid[v] = struct{}{}
	}

	valids := append(valid, dflt)
	usagestr := fmt.Sprintf(`%s ["%s"]`, usage, strings.Join(valids, `","`))

	flags.StringVar(&os.value, name, dflt, usagestr)

	return os
}

func (os optset) Value() (string, error) {
	if _, ok := os.valid[os.value]; !ok {
		return "", fmt.Errorf("invalid selection for %s: %s", os.name, os.value)
	}
	return os.value, nil
}
