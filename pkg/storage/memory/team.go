package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sync"
)

// teamRepository implements domain.TeamRepository interface
type teamRepository struct {
	teams map[int]domain.Team
	sync.RWMutex
}

// NewTeamRepository returns new instance of domain.TeamRepository
func NewTeamRepository() domain.TeamRepository {
	return &teamRepository{
		teams: make(map[int]domain.Team),
	}
}

// Add saves given team into memory storage, or returns an error on failure
func (tr *teamRepository) Add(team domain.Team) error {
	tr.Lock()
	defer tr.Unlock()

	if _, ok := tr.teams[team.ID]; ok {
		return storage.ErrTeamAlreadyExists
	}
	tr.teams[team.ID] = team
	return nil
}

// Update updates team with matching ID in memory storage, of return an error on failure
func (tr *teamRepository) Update(team domain.Team) error {
	tr.Lock()
	defer tr.Unlock()

	if _, ok := tr.teams[team.ID]; ok {
		tr.teams[team.ID] = team
		return nil
	}

	return storage.ErrTeamNotFound
}

// GetByID returns team with matching ID, or returns an error on failure
func (tr *teamRepository) GetByID(ID int) (domain.Team, error) {
	tr.RLock()
	defer tr.RUnlock()

	if team, ok := tr.teams[ID]; ok {
		return team, nil
	}

	return domain.Team{}, storage.ErrTeamNotFound
}

// GetCount returns number of team records in memory storage
func (tr *teamRepository) GetCount() (int, error) {
	tr.RLock()
	defer tr.RUnlock()

	return len(tr.teams), nil
}
