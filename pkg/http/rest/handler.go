package rest

import (
	"bytes"
	"encoding/json"
	"fpl-live-tracker/pkg/services/manager"
	"fpl-live-tracker/pkg/services/player"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func Handler(ps player.PlayerService, ms manager.ManagerService) http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", Homepage()).Methods("GET")
	r.HandleFunc("/api/players", GetPlayers(ps)).Methods("GET")
	r.HandleFunc("/api/manager/{id:[0-9]+}", GetManager(ms)).Methods("GET")
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

		players, err := ps.GetAll()
		if err != nil {
			http.Error(w, "failed to get players", http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(players)
		if err != nil {
			http.Error(w, "failed to marshal data", http.StatusInternalServerError)
			return
		}
		io.Copy(w, bytes.NewReader(j))
	}
}

//
func GetManager(ms manager.ManagerService) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		vars := mux.Vars(r)
		managerIDstr := vars["id"]
		managerID, err := strconv.ParseUint(managerIDstr, 10, 32)
		if err != nil {
			http.Error(w, "failed to parse url", http.StatusInternalServerError)
			return
		}

		m, err := ms.GetByID(int(managerID))
		if err != nil {
			http.Error(w, "failed to get manager", http.StatusInternalServerError)
			return
		}

		j, err := json.Marshal(m)
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
