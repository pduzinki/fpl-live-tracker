package domain

import "time"

type Gameweek struct {
	ID           int
	Name         string
	Finished     bool
	IsCurrent    bool
	IsNext       bool
	DeadlineTime time.Time
}
