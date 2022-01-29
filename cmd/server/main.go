package main

import (
	"log"
	"net/http"

	"fpl-live-tracker/pkg/http/rest"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/services/fixture"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/services/tracker"
	"fpl-live-tracker/pkg/storage/memory"
	"fpl-live-tracker/pkg/wrapper"
)

func main() {
	log.Println("fpl-live-tracker started")

	w := wrapper.NewWrapper()

	cr := memory.NewClubRepository()
	cs, err := club.NewClubService(cr, w)
	if err != nil {
		log.Fatalln("error: failed to init club service")
	}

	fr := memory.NewFixtureRepository()
	fs, err := fixture.NewFixtureService(fr, cs, w)
	if err != nil {
		log.Fatalln("error: failed to init fixture service")
	}

	gs, err := gameweek.NewGameweekService(w)
	if err != nil {
		log.Fatalln("error: failed to init gameweek service")
	}

	pr := memory.NewPlayerRepository()
	ps, err := player.NewPlayerService(w, pr, cs, fs, gs)
	if err != nil {
		log.Fatalln("error: failed to init player service")
	}

	tracker, err := tracker.NewTracker(
		tracker.WithPlayerService(ps),
		tracker.WithClubService(cs),
		tracker.WithFixtureService(fs),
		tracker.WithGameweekService(gs))
	if err != nil {
		log.Fatalf("failed to init tracker: %v\n", err)
	}
	go tracker.Track()

	router := rest.Handler()

	log.Println("fpl-live-tracker now listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
