package domain

import (
	"fmt"
	"time"
)

// Fixture represents a single match of football between two clubs in the Premier League
type Fixture struct {
	ID    int
	Info  FixtureInfo
	Stats map[string]FixtureStat
}

type FixtureRepository interface {
	Add(fixture Fixture) error
	AddMany(fixtures []Fixture) error
	Update(fixture Fixture) error
	GetByGameweek(gameweekID int) ([]Fixture, error)
	GetByID(fixtureID int) (Fixture, error)
}

//
type FixtureInfo struct {
	GameweekID          int
	ClubHome            Club
	ClubAway            Club
	Started             bool
	FinishedProvisional bool
	Finished            bool
	KickoffTime         time.Time
}

// FixtureStat represents particular fixture statistic (e.g. goals scored, assists, or bonus points)
type FixtureStat struct {
	Name             string
	HomePlayersStats []FixtureStatPair
	AwayPlayersStats []FixtureStatPair
}

// FixtureStatPair represents particular instance of fixture statistic, and player responsible for it (e.g. number of goals scored by Harry Kane)
type FixtureStatPair struct {
	PlayerID int
	Value    int
}

func (f Fixture) String() string {
	return fmt.Sprintf("{%v}", f.Info)
}
