package gameweek

import (
	"errors"
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/wrapper"
)

var ErrGameweekAllFinished error = errors.New("gameweek: all gameweeks finished") // game finished
var ErrGameweekNoneOngoing error = errors.New("gameweek: none ongoing gameweeks")

type GameweekService interface {
	GetOngoingGameweek() (domain.Gameweek, error)
	GetNextGameweek() (domain.Gameweek, error)
}

type gameweekService struct {
	wrapper wrapper.Wrapper
}

func NewGameweekService(w wrapper.Wrapper) GameweekService {
	return &gameweekService{
		wrapper: w,
	}
}

// GetOngoingGameweek returns current, ongoing gameweek.
func (gs *gameweekService) GetOngoingGameweek() (domain.Gameweek, error) {
	gameweeks, err := gs.wrapper.GetGameweeks()
	if err != nil {
		return domain.Gameweek{}, err // propagate error
	}

	for _, gw := range gameweeks {
		if gw.IsCurrent && !gw.Finished {
			return gw, nil
		}
	}

	return domain.Gameweek{}, ErrGameweekNoneOngoing
}

// GetNextGameweek returns subsequent gameweek
func (gs *gameweekService) GetNextGameweek() (domain.Gameweek, error) {
	gameweeks, err := gs.wrapper.GetGameweeks()
	if err != nil {
		return domain.Gameweek{}, err // propagate error
	}

	for _, gw := range gameweeks {
		if gw.IsNext {
			return gw, nil
		}
	}

	return domain.Gameweek{}, ErrGameweekAllFinished
}
