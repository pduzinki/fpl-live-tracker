package domain

// Fixture represents a single match of football between two clubs in the Premier League
type Fixture struct {
	GameweekID int
	ID         int
	ClubAway   Club
	ClubHome   Club
	// Started    bool
	// Finished   bool
}

type FixtureRepository interface {
	Add(fixture Fixture) error
	AddMany(fixtures []Fixture) error
	// Update() error // TODO add later
	GetByGameweek(gameweekID int) ([]Fixture, error)
	// TODO maybe add Get(s string) (Fixture, error) // s in form LEILIV, TOTCHE
}
