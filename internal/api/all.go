package api

import (
	"encoding/json"
	"github.com/frennkie/blitzd/internal/data"
	"github.com/frennkie/blitzd/internal/metric"
	"net/http"
)

func All() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		all := struct {
			data.Lnd
			data.Network
			data.System
		}{metric.Lnd, metric.Network, metric.System}

		js, err := json.Marshal(all)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(js)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
