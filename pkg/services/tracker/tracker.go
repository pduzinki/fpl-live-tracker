package tracker

import (
	"errors"
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/services/gameweek"
	"log"
	"time"
)

type TrackerConfigFunc func(t *Tracker) error

type Tracker struct {
	GwService gameweek.GameweekService
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
			return errors.New("tracker init error: gwService is nil")
		}
		t.GwService = gwService
		return nil
	}
}

// Track is responsible for keeping all the data from FPL up-to-date.
// Should be run as a goroutine.
func (t *Tracker) Track() {
	log.Println("hello from Track()")

	// check if there is current, ongoing gameweek, if yes, then proceed with processing data from it
	// if there is no ongoing gameweek, check when the next gameweek starts, and sleep until then
	// if there is no next gameweek, this means game ended

	var gameweek domain.Gameweek

	gameweek, err := t.GwService.GetOngoingGameweek()
	if err != nil {
		log.Println(err) // TODO handle err properly

		gameweek, err = t.GwService.GetNextGameweek()
		if err != nil {
			log.Println(err) // TODO handle err properly
		}
	}

	now := time.Now()
	if now.Before(gameweek.DeadlineTime) {
		diff := gameweek.DeadlineTime.Sub(now)
		log.Printf("Next gameweek starts in %v", diff)
	}

	// get current gameweek fixtures
	// tbd next steps

	log.Println("goodbye from Track()")
}
