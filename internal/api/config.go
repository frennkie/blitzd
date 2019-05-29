package api

import (
	"encoding/json"
	"github.com/spf13/viper"
	"net/http"
)

func Config() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		js, err := json.Marshal(viper.AllSettings())
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
