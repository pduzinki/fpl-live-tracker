package player

import (
	"fpl-live-tracker/pkg/wrapper"

	domain "fpl-live-tracker/pkg"
)

type PlayerService interface {
	Update() error
	UpdateGameweekData() error
	GetByID(ID int) (domain.Player, error)
}

type playerService struct {
	wrapper wrapper.Wrapper
	pr      domain.PlayerRepository
}

//
func NewPlayerService(w wrapper.Wrapper, pr domain.PlayerRepository) PlayerService {
	return &playerService{
		wrapper: w,
		pr:      pr,
	}
}

//
func (ps *playerService) Update() error {
	// TODO
	return nil
}

//
func (ps *playerService) UpdateGameweekData() error {
	// TODO
	return nil
}

//
func (ps *playerService) GetByID(ID int) (domain.Player, error) {
	// TODO
	return domain.Player{}, nil
}
