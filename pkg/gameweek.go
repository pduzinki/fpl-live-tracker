package tracker

type Gameweek struct {
	ID           int
	Name         string
	Finished     bool
	IsCurrent    bool
	IsNext       bool
	DeadlineTime string // TODO change to time.Time
}
