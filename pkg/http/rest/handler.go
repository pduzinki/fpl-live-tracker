package rest

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", Homepage()).Methods("GET")
	r.HandleFunc("/api/manager/{id:[0-9]+}", GetManager()).Methods("GET")
	r.HandleFunc("/api/league/{id:[0-9]+}", GetLeague()).Methods("GET")

	return r
}

//
func Homepage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		s := "fpl-live-tracker test"

		j, err := json.Marshal(s)
		if err != nil {
			http.Error(w, "failed to marshal data", http.StatusInternalServerError)
			return
		}

		io.Copy(w, bytes.NewReader(j))
	}
}

//
func GetManager() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		s := "fpl-live-tracker getmanager test"

		j, err := json.Marshal(s)
		if err != nil {
			http.Error(w, "failed to marshal data", http.StatusInternalServerError)
			return
		}

		io.Copy(w, bytes.NewReader(j))
	}
}

//
func GetLeague() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		s := "fpl-live-tracker getleague test"

		j, err := json.Marshal(s)
		if err != nil {
			http.Error(w, "failed to marshal data", http.StatusInternalServerError)
			return
		}

		io.Copy(w, bytes.NewReader(j))
	}
}
