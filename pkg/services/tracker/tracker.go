package tracker

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
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
	var gw domain.Gameweek
	var timeToUpdateManagers bool

	for {
		err := t.Gs.Update()
		if err != nil {
			log.Println("tracker service: failed to update gameweek data:", err)
		}

		currentGw, err := t.Gs.GetCurrentGameweek()
		if err != nil {
			log.Println("tracker service: failed to update gameweek data:", err)
		}

		if gw != currentGw || time.Now().Before(currentGw.DeadlineTime) {
			gw = currentGw
			timeToUpdateManagers = true
		}

		if currentGw.Finished {
			// after gameweek is finished, it's time to update manager's information
			err = t.Ms.UpdateInfos()
			if err != nil {
				log.Println("tracker service: failed to update manager's information:", err)
			}

			nextGw, err := t.Gs.GetNextGameweek()
			if err == gameweek.ErrNoNextGameweek {
				log.Println("tracker service: game ended, tracking stopped")
				return
			} else if err != nil {
				log.Println("tracker service: failed to get next gameweek:", err)
				continue
			}

			time.Sleep(time.Until(nextGw.DeadlineTime))
			continue
		}

		err = t.Fs.Update()
		if err != nil {
			log.Println("tracker service: failed to update fixture data:", err)
		}

		err = t.Ps.UpdateInfos()
		if err != nil {
			log.Println("tracker service: failed to update player data:", err)
		}

		err = t.Ps.UpdateStats()
		if err != nil {
			log.Println("tracker service: failed to update player data:", err)
		}

		if timeToUpdateManagers {
			err = t.Ms.AddNew()
			if err != nil {
				log.Println("tracker service: failed to add new manager's data:", err)
			}

			err = t.Ms.UpdateTeams()
			if err != nil {
				log.Println("tracker service: failed to update manager's teams:", err)
			}

			log.Println("tracker service: manager's teams updated")
			timeToUpdateManagers = false
		}

		liveFixtures, err := t.Fs.GetLiveFixtures(currentGw.ID)
		if err != nil {
			log.Println("tracker service: failed to get live fixtures:", err)
		}

		if len(liveFixtures) > 0 {
			err = t.Ms.UpdatePoints()
			if err != nil {
				log.Println("tracker service: failed to update manager's points:", err)
			}
		}

		log.Println("tracker service: FPL API data updated")
		time.Sleep(1 * time.Minute)
	}
}
