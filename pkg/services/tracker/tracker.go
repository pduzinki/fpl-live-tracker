package tracker

import (
	"errors"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/services/fixture"
	"fpl-live-tracker/pkg/services/gameweek"
	"log"
)

type TrackerConfigFunc func(t *Tracker) error

type Tracker struct {
	Cs club.ClubService
	Fs fixture.FixtureService
	Gs gameweek.GameweekService
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

//
func WithClubService(cs club.ClubService) TrackerConfigFunc {
	return func(t *Tracker) error {
		if cs == nil {
			return errors.New("tracker init error: Club Service is nil")
		}
		t.Cs = cs
		return nil
	}
}

//
func WithFixtureService(fs fixture.FixtureService) TrackerConfigFunc {
	return func(t *Tracker) error {
		if fs == nil {
			return errors.New("tracker init error: Fixture Service is nil")
		}
		t.Fs = fs
		return nil
	}
}

//
func WithGameweekService(gwService gameweek.GameweekService) TrackerConfigFunc {
	return func(t *Tracker) error {
		if gwService == nil {
			return errors.New("tracker init error: Gameweek Service is nil")
		}
		t.Gs = gwService
		return nil
	}
}

// Track is responsible for keeping all the data from FPL up-to-date.
// Should be run as a goroutine.
func (t *Tracker) Track() {
	log.Println("hello from Track()")

	err := t.Gs.Update()
	if err != nil {
		log.Println("tracker service: failed to update gameweek data", err)
	}

	err = t.Fs.Update()
	if err != nil {
		log.Println("tracker service: failed to update gameweek data")
	}

	gw, err := t.Gs.GetCurrentGameweek()
	if err != nil {
		log.Println("tracker service: failed to get current gameweek", err)
	}
	log.Println("current gameweek:", gw)

	fixtures, err := t.Fs.GetFixturesByGameweek(gw.ID)
	if err != nil {
		log.Println(err)
	}
	for _, f := range fixtures {
		log.Println(f)
	}

	log.Println("goodbye from Track()")
}
