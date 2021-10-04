package main

import (
	"fmt"
	"regexp"
	"strings"
)

const (
	labelNameRegex = "[a-zA-Z_][a-zA-Z0-9_]*"
)

func processMap(s string) (map[string]string, error) {
	if len(s) == 0 {
		return nil, nil
	}

	m := make(map[string]string)

	for _, pair := range strings.Split(s, ",") {
		kvpair := strings.Split(pair, ":")
		if len(kvpair) != 2 {
			return nil, fmt.Errorf("invalid map item: %q", pair)
		}
		name, value := kvpair[0], kvpair[1]

		ok, err := testLabelName(name)
		if err != nil {
			return nil, fmt.Errorf("could not check label name: %v", err)
		}
		if !ok {
			return nil, fmt.Errorf("name %q does not match regex %q", name, labelNameRegex)
		}

		m[name] = value
	}
	return m, nil
}

func testLabelName(n string) (bool, error) {
	return regexp.MatchString("^"+labelNameRegex+"$", n)
}
