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
}

//
func NewPlayerService(w wrapper.Wrapper) PlayerService {
	return &playerService{
		wrapper: w,
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
