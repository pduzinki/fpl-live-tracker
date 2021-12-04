package gameweek

import "fpl-live-tracker/pkg/wrapper"

type GameweekService interface {
	GetCurrentGameweek()
	GetNextGameweek()
}

type gameweekService struct {
	wrapper wrapper.Wrapper
}

func NewGameweekService(w wrapper.Wrapper) GameweekService {
	return &gameweekService{
		wrapper: w,
	}
}

// GetCurrentGameweek returns current, ongoing gameweek.
func (gs *gameweekService) GetCurrentGameweek() {

}

// GetNextGameweek returns subsequent gameweek
func (gs *gameweekService) GetNextGameweek() {

}
