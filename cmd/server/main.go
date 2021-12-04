package main

import (
	"fmt"
	"log"
	"net/http"

	"fpl-live-tracker/pkg/http/rest"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/tracker"
	"fpl-live-tracker/pkg/storage/memory"
	"fpl-live-tracker/pkg/wrapper"
)

func main() {
	fmt.Println("fpl-live-tracker started")

	mr, err := memory.NewManagerRepository()
	if err != nil {
		panic(err)
	}
	_ = mr

	wrapper := wrapper.NewWrapper(wrapper.DefaultURL)

	gwService := gameweek.NewGameweekService(wrapper)

	tracker, err := tracker.NewTracker(tracker.WithGameweekService(gwService))
	if err != nil {
		log.Fatalf("failed to init tracker: %v\n", err)
	}
	go tracker.Track()

	router := rest.Handler()

	log.Println("fpl-live-tracker now listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
