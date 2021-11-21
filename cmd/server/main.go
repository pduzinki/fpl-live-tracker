package main

import (
	"fmt"
	tracker "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/http/rest"
	"fpl-live-tracker/pkg/storage/memory"
	"log"
	"net/http"
)

func main() {
	fmt.Println("fpl-live-tracker started")

	mr, err := memory.NewManagerRepository()
	if err != nil {
		panic(err)
	}
	_ = mr

	tracker, err := tracker.NewTracker()
	if err != nil {
		log.Fatalf("failed to init tracker: %v\n", err)
	}
	go tracker.Track()

	router := rest.Handler()

	log.Println("fpl-live-tracker now listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
