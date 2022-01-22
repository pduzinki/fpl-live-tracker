package domain

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
	Stats    Stats
}

type PlayerRepository interface {
	Add(player Player) error
	Update(player Player) error
	UpdateStats(playerID int, stats Stats) error
	GetByID(ID int) (Player, error)
}

// Stats contains data about Player's performance during Gameweek
type Stats struct {
	// Finished    bool
	Minutes     int
	TotalPoints int // doesn't include bonus points just yet
}
