package rest

import (
	"bytes"
	"encoding/json"
	"fpl-live-tracker/pkg/services/player"
	"io"
	"net/http"

	"github.com/gorilla/mux"
)

func Handler(ps player.PlayerService) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", Homepage()).Methods("GET")
	r.HandleFunc("/api/players", GetPlayers(ps)).Methods("GET")
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

func GetPlayers(ps player.PlayerService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		players, _ := ps.GetAll()

		j, err := json.Marshal(players)
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
