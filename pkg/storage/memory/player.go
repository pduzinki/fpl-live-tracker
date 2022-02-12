package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sort"
	"sync"
)

// playerRepository implements domain.PlayerRepository interface
type playerRepository struct {
	players map[int]domain.Player
	sync.Mutex
}

// NewPlayerRepository returns new instance of domain.PlayerRepository
func NewPlayerRepository() domain.PlayerRepository {
	return &playerRepository{
		players: make(map[int]domain.Player),
	}
}

// Add saves given player into memory storage, returns error on failure
func (pr *playerRepository) Add(player domain.Player) error {
	if _, ok := pr.players[player.ID]; ok {
		return storage.ErrPlayerAlreadyExists
	}

	pr.Lock()
	pr.players[player.ID] = player
	pr.Unlock()

	return nil
}

// Update updates player with matching ID in memory storage, or returns error on failure
func (pr *playerRepository) UpdateInfo(playerID int, info domain.PlayerInfo) error {
	pr.Lock()
	defer pr.Unlock()

	if player, ok := pr.players[playerID]; ok {
		player.Info = info
		pr.players[playerID] = player
		return nil
	}

	return storage.ErrPlayerNotFound
}

// UpdateStats updates stats of player with given playerID, or returns error on failure
func (pr *playerRepository) UpdateStats(playerID int, stats domain.PlayerStats) error {
	if player, ok := pr.players[playerID]; ok {
		player.Stats = stats
		pr.Lock()
		pr.players[playerID] = player
		pr.Unlock()
		return nil
	}

	return storage.ErrPlayerNotFound
}

// GetByID returns player with given ID, or returns error on failure
func (pr *playerRepository) GetByID(ID int) (domain.Player, error) {
	pr.Lock()
	defer pr.Unlock()
	if player, ok := pr.players[ID]; ok {
		return player, nil
	}

	return domain.Player{}, storage.ErrPlayerNotFound
}

// GetAll returns slice of all players, or error on failure
func (pr *playerRepository) GetAll() ([]domain.Player, error) {
	players := make([]domain.Player, 0, len(pr.players))

	pr.Lock()
	for _, p := range pr.players {
		players = append(players, p)
	}
	pr.Unlock()

	sort.Slice(players, func(i, j int) bool {
		return players[i].ID < players[j].ID
	})

	return players, nil
}
