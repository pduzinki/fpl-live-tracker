package mock

import "fpl-live-tracker/pkg/domain"

type PlayerRepository struct {
	AddFn         func(player domain.Player) error
	UpdateInfoFn  func(playerID int, info domain.PlayerInfo) error
	UpdateStatsFn func(playerID int, stats domain.PlayerStats) error
	GetByIDFn     func(ID int) (domain.Player, error)
	GetAllFn      func() ([]domain.Player, error)
}

func (pr *PlayerRepository) Add(player domain.Player) error {
	return pr.AddFn(player)
}

func (pr *PlayerRepository) UpdateInfo(playerID int, info domain.PlayerInfo) error {
	return pr.UpdateInfoFn(playerID, info)
}

func (pr *PlayerRepository) UpdateStats(playerID int, stats domain.PlayerStats) error {
	return pr.UpdateStatsFn(playerID, stats)
}

func (pr *PlayerRepository) GetByID(ID int) (domain.Player, error) {
	return pr.GetByIDFn(ID)
}

func (pr *PlayerRepository) GetAll() ([]domain.Player, error) {
	return pr.GetAllFn()
}
