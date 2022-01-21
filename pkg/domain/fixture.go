package domain

import (
	"fmt"
	"time"
)

// Fixture represents a single match of football between two clubs in the Premier League
type Fixture struct {
	GameweekID          int
	ID                  int
	ClubHome            Club
	ClubAway            Club
	Started             bool
	Finished            bool
	FinishedProvisional bool
	KickoffTime         time.Time
	Stats               []FixtureStat // TODO this can be a map
}

type FixtureRepository interface {
	Add(fixture Fixture) error
	AddMany(fixtures []Fixture) error
	Update(fixture Fixture) error
	GetByGameweek(gameweekID int) ([]Fixture, error)
	// TODO maybe add Get(s string) (Fixture, error) // s in form LEILIV, TOTCHE
}

//
type FixtureStat struct {
	Name             string
	HomePlayersStats []FixtureStatValue
	AwayPlayersStats []FixtureStatValue
}

//
type FixtureStatValue struct {
	PlayerID int
	Value    int
}

func (f Fixture) String() string {
	return fmt.Sprintf("{%d %d %s %s %t %t %t}",
		f.GameweekID,
		f.ID,
		f.ClubHome.Shortname,
		f.ClubAway.Shortname,
		f.Started,
		f.FinishedProvisional,
		f.Finished)
}
