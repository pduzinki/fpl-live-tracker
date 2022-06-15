package main

import (
	"log"
	"net/http"

	"fpl-live-tracker/pkg/config"
	"fpl-live-tracker/pkg/http/rest"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/services/fixture"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/manager"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/services/team"
	"fpl-live-tracker/pkg/services/tracker"
	"fpl-live-tracker/pkg/storage/memory"
	"fpl-live-tracker/pkg/storage/mongo"
	"fpl-live-tracker/pkg/wrapper"
)

func main() {
	log.Println("fpl-live-tracker started")

	cfg := config.Load()
	_ = cfg

	wr := wrapper.NewWrapper()

	cr := memory.NewClubRepository()
	cs, err := club.NewClubService(cr, wr)
	if err != nil {
		log.Fatalln("error: failed to init club service")
	}

	fr := memory.NewFixtureRepository()
	fs, err := fixture.NewFixtureService(fr, cs, wr)
	if err != nil {
		log.Fatalln("error: failed to init fixture service")
	}

	gs, err := gameweek.NewGameweekService(wr)
	if err != nil {
		log.Fatalln("error: failed to init gameweek service")
	}

	pr := memory.NewPlayerRepository()
	ps, err := player.NewPlayerService(wr, pr, cs, fs, gs)
	if err != nil {
		log.Fatalln("error: failed to init player service")
	}

	mr, err := mongo.NewManagerRepository(cfg.MongoConfig)
	if err != nil {
		log.Fatalln("error: failed to init manager repository")
	}
	// mr := memory.NewManagerRepository()

	ms, err := manager.NewManagerService(mr, gs, wr)
	if err != nil {
		log.Fatalln("error: failed to init manager service")
	}

	tr, err := mongo.NewTeamRepository(cfg.MongoConfig)
	if err != nil {
		log.Fatalln("error: failed to init team repository")
	}

	ts := team.NewTeamService(tr, gs, ps, wr)

	tracker, err := tracker.NewTracker(
		tracker.WithPlayerService(ps),
		tracker.WithClubService(cs),
		tracker.WithFixtureService(fs),
		tracker.WithGameweekService(gs),
		tracker.WithManagerService(ms),
		tracker.WithTeamService(ts))
	if err != nil {
		log.Fatalf("failed to init tracker: %v\n", err)
	}
	go tracker.Track()

	router := rest.Handler(ps, ms, ts)

	log.Println("fpl-live-tracker now listening on port 3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}
