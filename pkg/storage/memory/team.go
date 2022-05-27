package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sync"
)

//
type teamRepository struct {
	teams map[int]domain.Team
	sync.RWMutex
}

//
func NewTeamRepository() domain.TeamRepository {
	return &teamRepository{
		teams: make(map[int]domain.Team),
	}
}

//
func (tr *teamRepository) Add(team domain.Team) error {
	tr.Lock()
	defer tr.Unlock()

	if _, ok := tr.teams[team.ID]; ok {
		return storage.ErrTeamAlreadyExists
	}
	tr.teams[team.ID] = team
	return nil
}

//
func (tr *teamRepository) Update(ID int, team domain.Team) error {
	tr.Lock()
	defer tr.Unlock()

	if _, ok := tr.teams[ID]; ok {
		tr.teams[ID] = team
		return nil
	}

	return storage.ErrTeamNotFound
}

//
func (tr *teamRepository) GetByID(ID int) (domain.Team, error) {
	tr.RLock()
	defer tr.RUnlock()

	if team, ok := tr.teams[ID]; ok {
		return team, nil
	}

	return domain.Team{}, storage.ErrTeamNotFound
}

//
func (tr *teamRepository) GetCount() (int, error) {
	tr.RLock()
	defer tr.RUnlock()

	return len(tr.teams), nil
}
