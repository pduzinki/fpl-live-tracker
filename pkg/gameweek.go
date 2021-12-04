package domain

type Gameweek struct {
	ID           int
	Name         string
	Finished     bool
	IsCurrent    bool
	IsNext       bool
	DeadlineTime string // TODO change to time.Time
}

// TODO move domain types to separate folder
