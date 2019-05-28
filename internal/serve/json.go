package serve

import (
	"encoding/json"
	"github.com/frennkie/blitzinfod/internal/data"
	"net/http"
)

func StaticApi(metrics *data.Cache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		js, err := json.Marshal(metrics)
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
