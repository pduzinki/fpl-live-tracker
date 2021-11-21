package tracker

import (
	"log"
	"time"
)

type TrackerConfigFunc func(t *Tracker) error

type Tracker struct {
	GwService GameweekService
}

//
func NewTracker(fns ...TrackerConfigFunc) (*Tracker, error) {
	t := Tracker{
		// GwService: GameweekServdice{},
	}

	for _, f := range fns {
		err := f(&t)
		if err != nil {
			return nil, err
		}
	}

	return &t, nil
}

// Track is responsible for keeping all the data from FPL up-to-date.
// Should be run as a goroutine.
func (t *Tracker) Track() {
	log.Println("hello from Track()")

	now := time.Now()
	_ = now

	log.Println("goodbye from Track()")
}
