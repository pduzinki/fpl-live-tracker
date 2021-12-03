package gameweek

type GameweekService interface {
	GetCurrentGameweek()
	GetNextGameweek()
}

type gameweekService struct {
}

func NewGameweekService() GameweekService {
	return &gameweekService{}
}

// GetCurrentGameweek returns current, ongoing gameweek.
func (gs *gameweekService) GetCurrentGameweek() {

}

// GetNextGameweek returns subsequent gameweek
func (gs *gameweekService) GetNextGameweek() {

}
