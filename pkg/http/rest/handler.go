package rest

import (
	"net/http"

	"github.com/gorilla/mux"
)

func Handler() http.Handler {
	r := mux.NewRouter()

	r.HandleFunc("/", nil).Methods("GET")

	return r
}

//
func Homepage() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

	}
}
