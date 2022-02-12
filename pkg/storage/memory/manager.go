package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sync"
)

// managerRepository implements domain.ManagerRepository interface
type managerRepository struct {
	managers map[int]domain.Manager
	sync.Mutex
}

// NewManagerRepository returns new instance of domain.ManagerRepository
func NewManagerRepository() domain.ManagerRepository {
	mr := managerRepository{
		managers: make(map[int]domain.Manager),
	}

	return &mr
}

// Add saves given manager into memory storage, or returns error on failure
func (mr *managerRepository) Add(manager domain.Manager) error {
	if _, ok := mr.managers[manager.ID]; ok {
		return storage.ErrManagerAlreadyExists
	}

	mr.Lock()
	mr.managers[manager.ID] = manager
	mr.Unlock()

	return nil
}

// AddMany saves all given managers into memory storage, or returns error on failure
func (mr *managerRepository) AddMany(managers []domain.Manager) error {
	for _, manager := range managers {
		err := mr.Add(manager)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update updates manager with matching ID in memory storage, or returns error on failure
func (mr *managerRepository) UpdateInfo(managerID int, info domain.ManagerInfo) error {
	mr.Lock()
	defer mr.Unlock()
	if m, ok := mr.managers[managerID]; ok {
		m.Info = info
		mr.managers[managerID] = m
		return nil
	}

	return storage.ErrManagerNotFound
}

// UpdateTeam updates team of manager with given ID, or returns error on failure
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

// GetById returns manager with given ID, or returns error on failure
func (mr *managerRepository) GetByID(id int) (domain.Manager, error) {
	if manager, ok := mr.managers[id]; ok {
		return manager, nil
	}

	return domain.Manager{}, storage.ErrManagerNotFound
}
