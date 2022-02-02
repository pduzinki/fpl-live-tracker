package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sync"
)

//
type managerRepository struct {
	managers map[int]domain.Manager
	sync.Mutex
}

// NewManagerRepository creates new repository for managers
func NewManagerRepository() domain.ManagerRepository {
	mr := managerRepository{
		managers: make(map[int]domain.Manager),
	}

	return &mr
}

//
func (mr *managerRepository) Add(manager domain.Manager) error {
	if _, ok := mr.managers[manager.ID]; ok {
		return storage.ErrManagerAlreadyExists
	}

	mr.Lock()
	mr.managers[manager.ID] = manager
	mr.Unlock()

	return nil
}

//
func (mr *managerRepository) AddMany(managers []domain.Manager) error {
	for _, manager := range managers {
		err := mr.Add(manager)
		if err != nil {
			return err
		}
	}
	return nil
}

//
func (mr *managerRepository) Update(manager domain.Manager) error {
	mr.Lock()
	defer mr.Unlock()
	if m, ok := mr.managers[manager.ID]; ok {
		m.Name = manager.Name
		m.TeamName = manager.TeamName
		mr.managers[manager.ID] = m
		return nil
	}

	return storage.ErrManagerNotFound
}

//
func (mr *managerRepository) UpdateTeam(managerID int, team domain.Team) error {
	mr.Lock()
	defer mr.Unlock()
	if m, ok := mr.managers[managerID]; ok {
		m.Team = team
		mr.managers[managerID] = m
		return nil
	}

	return storage.ErrManagerNotFound
}

//
func (mr *managerRepository) GetByID(id int) (domain.Manager, error) {
	if manager, ok := mr.managers[id]; ok {
		return manager, nil
	}

	return domain.Manager{}, storage.ErrManagerNotFound
}
