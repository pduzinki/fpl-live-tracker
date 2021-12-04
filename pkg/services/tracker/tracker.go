package tracker

import (
	"errors"
	"fpl-live-tracker/pkg/services/gameweek"
	"log"
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

	// now := time.Now()
	// log.Println(now)

	// get current gameweek
	// get current gameweek fixtures
	// tbd next steps

	log.Println("goodbye from Track()")
}
