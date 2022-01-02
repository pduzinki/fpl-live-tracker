package tracker

import (
	"errors"
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/services/fixture"
	"fpl-live-tracker/pkg/services/gameweek"
	"log"
	"time"
)

type TrackerConfigFunc func(t *Tracker) error

type Tracker struct {
	Gs gameweek.GameweekService
	Fs fixture.FixtureService
}

//
func NewTracker(fns ...TrackerConfigFunc) (*Tracker, error) {
	t := Tracker{}

	for _, f := range fns {
		err := f(&t)
		if err != nil {
			return nil, err
		}
	}

	return &t, nil
}

func WithGameweekService(gwService gameweek.GameweekService) TrackerConfigFunc {
	return func(t *Tracker) error {
		if gwService == nil {
			return errors.New("tracker init error: Gameweek Service is nil")
		}
		t.Gs = gwService
		return nil
	}
}

func WithFixtureService(fs fixture.FixtureService) TrackerConfigFunc {
	return func(t *Tracker) error {
		if fs == nil {
			return errors.New("tracker init error: Fixture Service is nil")
		}
		t.Fs = fs
		return nil
	}
}

// Track is responsible for keeping all the data from FPL up-to-date.
// Should be run as a goroutine.
func (t *Tracker) Track() {
	log.Println("hello from Track()")

	// find ongoing or next gameweek, if there is none, the game finished
	var gameweek domain.Gameweek

	gameweek, err := t.Gs.GetOngoingGameweek()
	if err != nil {
		log.Println(err) // TODO handle err properly

		gameweek, err = t.Gs.GetNextGameweek()
		if err != nil {
			log.Println(err) // TODO handle err properly
		}
	}

	now := time.Now()
	if now.Before(gameweek.DeadlineTime) {
		diff := gameweek.DeadlineTime.Sub(now)
		log.Printf("Gameweek %d starts in %v", gameweek.ID, diff)
	} else {
		log.Printf("Gameweek %d is live!", gameweek.ID)
	}

	// update fixtures
	err = t.Fs.Update(gameweek.ID)
	if err != nil {
		panic(err)
	}

	// get current gameweek fixtures
	fixtures, err := t.Fs.GetFixturesByGameweek(gameweek.ID)
	if err != nil {
		log.Println(err) // TODO handle err properly
	}
	for _, f := range fixtures {
		log.Println(f)
	}

	log.Println("goodbye from Track()")
}
