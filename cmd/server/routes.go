package main

import (
	"net/http"
)

func (a *app) routes(withREST bool) {
	a.router.HandleFunc("/", a.handleSD())
	a.router.HandleFunc("/v1/httpsd", a.handleSD())
	a.router.HandleFunc("/version", handleVersion())
	if withREST {
		a.router.HandleFunc("/v1/announce", a.handleAnnounce()).Methods("POST")
	}
}

func (a *app) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	a.router.ServeHTTP(w, req)
}
