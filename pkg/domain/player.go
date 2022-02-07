package domain

// PlayerPosition maps FPL API position IDs to more readable string representation
var PlayerPosition = map[int]string{
	1: "GKP",
	2: "DEF",
	3: "MID",
	4: "FWD",
}

// Player represents a human being that plays in one of Premier League clubs (e.g. Harry Kane, Mohamed Salah)
type Player struct {
	ID       int
	Name     string
	Position string
	Club     Club
	Stats    PlayerStats
}

type PlayerRepository interface {
	Add(player Player) error
	Update(player Player) error
	UpdateStats(playerID int, stats PlayerStats) error
	GetByID(ID int) (Player, error)
	GetAll() ([]Player, error)
}

// PlayerStats contains data about Player's performance during Gameweek
type PlayerStats struct {
	FixturesInfo []FixtureInfo
	Minutes      int
	TotalPoints  int
}
