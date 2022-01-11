package domain

// Player represents a human being that plays in one of Premier League clubs (e.g. Harry Kane, Mohamed Salah)
type Player struct {
	ID       int
	Name     string
	Position string
	Club     Club
}

type PlayerRepository interface {
	Add(player Player) error
	Update(player Player) error
	GetByID(ID int) (Player, error)
}
