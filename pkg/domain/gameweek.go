package domain

import "time"

// Gameweek represents a series of (most often 10) fixtures, that usually happen in a span of one weekend
type Gameweek struct {
	ID           int
	Name         string
	Finished     bool
	DeadlineTime time.Time
}
