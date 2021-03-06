package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sync"
)

// managerRepository implements domain.ManagerRepository interface
type managerRepository struct {
	managers map[int]domain.Manager
	sync.RWMutex
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
	mr.Lock()
	defer mr.Unlock()

	if _, ok := mr.managers[manager.ID]; ok {
		return storage.ErrManagerAlreadyExists
	}
	mr.managers[manager.ID] = manager

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
func (mr *managerRepository) Update(manager domain.Manager) error {
	mr.Lock()
	defer mr.Unlock()

	if _, ok := mr.managers[manager.ID]; ok {
		mr.managers[manager.ID] = manager
		return nil
	}

	return storage.ErrManagerNotFound
}

// GetByID returns manager with given ID, or returns error on failure
func (mr *managerRepository) GetByID(id int) (domain.Manager, error) {
	mr.RLock()
	defer mr.RUnlock()

	if manager, ok := mr.managers[id]; ok {
		return manager, nil
	}

	return domain.Manager{}, storage.ErrManagerNotFound
}

// GetCount returns number of manager records in memory storage
func (mr *managerRepository) GetCount() (int, error) {
	mr.RLock()
	defer mr.RUnlock()

	return len(mr.managers), nil
}
