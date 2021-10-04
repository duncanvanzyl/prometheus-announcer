package main

import (
	"net/http"
)

func (a *app) handleSD() http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		rw.Header().Set("Content-Type", "application/json")

		b, err := a.cs.JSON()
		if err != nil {
			http.Error(rw, "http service discovery error", http.StatusInternalServerError)
			return
		}
		rw.Write(b)
	}
}
