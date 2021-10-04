package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	sd "github.com/duncanvanzyl/prometheus-announcer/servicediscovery"

	"github.com/hashicorp/go-hclog"
)

func (a *app) handleAnnounce() http.HandlerFunc {
	type Announcement struct {
		ID string `json:"id"`
		sd.Config
	}

	return func(w http.ResponseWriter, req *http.Request) {
		defer req.Body.Close()
		bs, err := ioutil.ReadAll(req.Body)
		if err != nil {
			logger.Error("could not read announce request", "body", hclog.Fmt("%q", bs), "error", err)
			http.Error(w, "could not process request", http.StatusBadRequest)
			return
		}

		an := &Announcement{}
		err = json.Unmarshal(bs, an)
		if err != nil {
			logger.Error("could not unmarshal announce request", "body", hclog.Fmt("%q", bs), "error", err)
			http.Error(w, "could not process request", http.StatusBadRequest)
			return
		}

		if an.ID == "" {
			logger.Error("request has empty id", "announcement", an)
			http.Error(w, "empty id", http.StatusBadRequest)
			return
		}

		logger.Debug("got rest announcement", "announcement", hclog.Fmt("%+v", an))
		a.cs.AddTarget(an.ID, an.Targets, an.Labels)
	}
}
