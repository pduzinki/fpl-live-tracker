package gameweek

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"time"
)

type GameweekService interface {
	Update() error
	GetCurrentGameweek() (domain.Gameweek, error)
	GetNextGameweek() (domain.Gameweek, error)
}

type gameweekService struct {
	CurrentGameweek domain.Gameweek
	NextGameweek    domain.Gameweek
	wrapper         wrapper.Wrapper
}

func NewGameweekService(w wrapper.Wrapper) GameweekService {
	return &gameweekService{
		wrapper: w,
	}
}

func (gs *gameweekService) Update() error {
	wrapperGameweeks, err := gs.wrapper.GetGameweeks()
	if err != nil {
		log.Println("gameweek service:", err)
		return err
	}

	for _, gw := range wrapperGameweeks {
		if gw.IsCurrent {
			deadlineTime, err := time.Parse(time.RFC3339, gw.DeadlineTime)
			if err != nil {
				log.Println("gameweek service:", err)
				return err
			}

			gs.CurrentGameweek = domain.Gameweek{
				ID:           gw.ID,
				Name:         gw.Name,
				Finished:     gw.Finished,
				DeadlineTime: deadlineTime,
			}
		}

		if gw.IsNext {
			deadlineTime, err := time.Parse(time.RFC3339, gw.DeadlineTime)
			if err != nil {
				log.Println("gameweek service:", err)
				return err
			}

			gs.NextGameweek = domain.Gameweek{
				ID:           gw.ID,
				Name:         gw.Name,
				Finished:     gw.Finished,
				DeadlineTime: deadlineTime,
			}
		}
	}

	return nil
}

func (gs *gameweekService) GetCurrentGameweek() (domain.Gameweek, error) {
	return gs.CurrentGameweek, nil
}

func (gs *gameweekService) GetNextGameweek() (domain.Gameweek, error) {
	return gs.NextGameweek, nil
}
