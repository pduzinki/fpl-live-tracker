package gameweek

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"time"
)

var ErrGameweekNotUpdated = errors.New("gameweek service: gameweek data hasn't been updated")
var ErrNoNextGameweek = errors.New("gameweek service: no noxt gameweek found")

type GameweekService interface {
	Update() error
	GetCurrentGameweek() (domain.Gameweek, error)
	GetNextGameweek() (domain.Gameweek, error)
}

type gameweekService struct {
	CurrentGameweek domain.Gameweek
	NextGameweek    domain.Gameweek
	wr              wrapper.Wrapper
	noNextGameweek  bool
}

func NewGameweekService(w wrapper.Wrapper) (GameweekService, error) {
	gs := gameweekService{
		wr:             w,
		noNextGameweek: false,
	}

	err := gs.Update()
	if err != nil {
		log.Println("gameweek service: failed to init data")
		return nil, err
	}

	return &gs, nil
}

func (gs *gameweekService) Update() error {
	wrapperGameweeks, err := gs.wr.GetGameweeks()
	if err != nil {
		log.Println("gameweek service:", err)
		return err
	}

	nextGameweekFound := false
	for _, gw := range wrapperGameweeks {
		if gw.IsCurrent {
			currentGameweek, err := gs.convertToDomainGameweek(gw)
			if err != nil {
				return err
			}
			gs.CurrentGameweek = currentGameweek
		}

		if gw.IsNext {
			nextGameweek, err := gs.convertToDomainGameweek(gw)
			if err != nil {
				return err
			}
			gs.NextGameweek = nextGameweek
			nextGameweekFound = true
		}
	}

	if !nextGameweekFound {
		gs.noNextGameweek = true
	}
	return nil
}

func (gs *gameweekService) GetCurrentGameweek() (domain.Gameweek, error) {
	return gs.CurrentGameweek, nil
}

func (gs *gameweekService) GetNextGameweek() (domain.Gameweek, error) {
	if gs.noNextGameweek {
		return domain.Gameweek{}, ErrNoNextGameweek
	}
	return gs.NextGameweek, nil
}

func (gs *gameweekService) convertToDomainGameweek(gw wrapper.Gameweek) (domain.Gameweek, error) {
	deadlineTime, err := time.Parse(time.RFC3339, gw.DeadlineTime)
	if err != nil {
		log.Println("gameweek service:", err)
		return domain.Gameweek{}, err
	}

	return domain.Gameweek{
		ID:           gw.ID,
		Name:         gw.Name,
		Finished:     gw.Finished,
		DeadlineTime: deadlineTime,
	}, nil
}
