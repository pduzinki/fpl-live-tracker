package domain

type Fixture struct {
	GameweekID int
	ID         int
	Started    bool
	Finished   bool
	TeamA      int // TODO add TeamAName
	TeamH      int // TODO add TeamHName
}

type FixtureRepository interface {
	GetByGameweek(gameweekID int) ([]Fixture, error)
	// TODO maybe add Get(s string) (Fixture, error) // s in form LEILIV, TOTCHE
	Add(Fixture) error
	AddMany([]Fixture) error
}
