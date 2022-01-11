package memory

import (
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/storage"
	"sync"
)

//
type playerRepository struct {
	players map[int]domain.Player
	sync.Mutex
}

//
func NewPlayerRepository() domain.PlayerRepository {
	return &playerRepository{
		players: make(map[int]domain.Player),
	}
}

func (pr *playerRepository) Add(player domain.Player) error {
	if _, ok := pr.players[player.ID]; ok {
		return storage.ErrPlayerAlreadyExists
	}

	pr.Lock()
	pr.players[player.ID] = player
	pr.Unlock()

	return nil
}

func (pr *playerRepository) Update(player domain.Player) error {
	if _, ok := pr.players[player.ID]; ok {
		pr.players[player.ID] = player
		return nil
	}

	return storage.ErrPlayerNotFound
}

func (pr *playerRepository) GetByID(ID int) (domain.Player, error) {
	if player, ok := pr.players[ID]; ok {
		return player, nil
	}

	return domain.Player{}, storage.ErrPlayerNotFound
}
