package main

import (
	"log"
	"net/http"

	"fpl-live-tracker/pkg/http/rest"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/services/fixture"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/tracker"
	"fpl-live-tracker/pkg/storage/memory"
	"fpl-live-tracker/pkg/wrapper"
)

func main() {
	log.Println("fpl-live-tracker started")

	w := wrapper.NewWrapper(wrapper.DefaultURL)

	cr := memory.NewClubRepository()
	cs, err := club.NewClubService(cr, w)
	if err != nil {
		panic(err)
	}
	_ = cs

	gwService := gameweek.NewGameweekService(w)
	fs := fixture.NewFixtureService(w)

	tracker, err := tracker.NewTracker(
		tracker.WithGameweekService(gwService),
		tracker.WithFixtureService(fs))
	if err != nil {
		log.Fatalf("failed to init tracker: %v\n", err)
	}
	go tracker.Track()

	router := rest.Handler()

	log.Println("fpl-live-tracker now listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
