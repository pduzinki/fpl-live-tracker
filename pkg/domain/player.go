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
	ID    int         `bson:"ID"`
	Info  PlayerInfo  `bson:"PlayerInfo"`
	Stats PlayerStats `bson:"PlayerStats"`
}

//
type PlayerInfo struct {
	Name     string `bson:"Name"`
	Position string `bson:"Position"`
	Club     Club   `bson:"Club"`
}

// PlayerStats contains data about Player's performance during Gameweek
type PlayerStats struct {
	FixturesInfo []FixtureInfo `bson:"FixtureInfo"`
	Minutes      int           `bson:"Minutes"`
	TotalPoints  int           `bson:"TotalPoints"`
}

type PlayerRepository interface {
	Add(player Player) error
	UpdateInfo(playerID int, info PlayerInfo) error
	UpdateStats(playerID int, stats PlayerStats) error
	GetByID(ID int) (Player, error)
	GetAll() ([]Player, error)
}
