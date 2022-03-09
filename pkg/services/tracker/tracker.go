package tracker

import (
	"errors"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/services/fixture"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/manager"
	"fpl-live-tracker/pkg/services/player"
	"log"
	"time"
)

type TrackerConfigFunc func(t *Tracker) error

type Tracker struct {
	Ps player.PlayerService
	Cs club.ClubService
	Fs fixture.FixtureService
	Gs gameweek.GameweekService
	Ms manager.ManagerService
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
func WithPlayerService(ps player.PlayerService) TrackerConfigFunc {
	return func(t *Tracker) error {
		if ps == nil {
			return errors.New("tracker init error: Player Service is nil")
		}
		t.Ps = ps
		return nil
	}
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

//
func WithManagerService(ms manager.ManagerService) TrackerConfigFunc {
	return func(t *Tracker) error {
		if ms == nil {
			return errors.New("tracker init error: Manager Service is nil")
		}
		t.Ms = ms
		return nil
	}
}

// Track is responsible for keeping all the data from FPL up-to-date.
// Should be run as a goroutine.
func (t *Tracker) Track() {
	/*
		schedule
		update gameweeks data on loop

		before gameweek starts:
			update managers info data on loop
		gameweek is live:
			update managers teams data once

			update fixtures data on loop

			if there are live fixtures:
				update players data on loop
				update managers points data on loop
	*/

	// initial updates, in case the app had to be restarted during the gameweek
	err := t.Gs.Update()
	if err != nil {
		log.Println("tracker:", err)
	}

	err = t.Fs.Update()
	if err != nil {
		log.Println("tracker:", err)
	}

	err = t.Ps.UpdateInfos()
	if err != nil {
		log.Println("tracker:", err)
	}

	err = t.Ps.UpdateStats()
	if err != nil {
		log.Println("tracker:", err)
	}

	err = t.Ms.UpdateInfos()
	if err != nil {
		log.Println("tracker:", err)
	}

	err = t.Ms.UpdateTeams()
	if err != nil {
		log.Println("tracker:", err)
	}

	err = t.Ms.UpdatePoints()
	if err != nil {
		log.Println("tracker:", err)
	}

	var timeToUpdateManagersInfos bool
	var timeToUpdateManagersTeams bool

	for {
		err = t.Gs.Update()
		if err != nil {
			log.Println("tracker:", err)
		}
		currentGw, err := t.Gs.GetCurrentGameweek()
		if err != nil {
			log.Println("tracker:", err)
		}

		if currentGw.Finished { // before gameweek starts / after gameweek is finished
			timeToUpdateManagersTeams = true
			if timeToUpdateManagersInfos {
				err = t.Ms.UpdateInfos() // once per gameweek
				if err != nil {
					log.Println("tracker:", err)
				}
				timeToUpdateManagersInfos = false
			}
			err = t.Ms.AddNew() // many times between gameweeks
			if err != nil {
				log.Println("tracker:", err)
			}
			err = t.Ps.UpdateInfos() // many times between gameweeks
			if err != nil {
				log.Println("tracker:", err)
			}

			time.Sleep(1 * time.Hour)
		} else { // gameweek is live
			timeToUpdateManagersInfos = true
			if timeToUpdateManagersTeams {
				err = t.Ms.UpdateTeams() // once per gameweek
				if err != nil {
					log.Println("tracker:", err)
				}
				timeToUpdateManagersTeams = false
			}
			err = t.Fs.Update() //many times between gameweeks
			if err != nil {
				log.Println("tracker:", err)
			}

			fixtures, err := t.Fs.GetLiveFixtures(currentGw.ID)
			if err != nil {
				log.Println("tracker:", err)
			}

			if len(fixtures) > 0 {
				err = t.Ps.UpdateStats()
				if err != nil {
					log.Println("tracker:", err)
				}
				err = t.Ms.UpdatePoints()
				if err != nil {
					log.Println("tracker:", err)
				}
				time.Sleep(1 * time.Minute)
				continue
			} else {
				time.Sleep(5 * time.Minute)
				continue
			}
		}
	}
}
