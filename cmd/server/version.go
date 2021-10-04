package main

import (
	"fmt"
	"net/http"
)

// version and build info
var (
	version   string // sha1 revision used to build the program
	buildTime string // when the executable was built
	gitCommit string // the hash of the current commit
)

func buildInfo() string {
	return fmt.Sprintf("Built on %s from %s\nGit Commit: %s\n", buildTime, version, gitCommit)
}

func handleVersion() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(rw, "%s", buildInfo())
	}
}
