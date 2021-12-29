package memory

import (
	domain "fpl-live-tracker/pkg"
	"sync"
)

// type Manager struct {
// 	FplID    int
// 	FullName string
// }

//
type managerRepository struct {
	managers map[int]domain.Manager
	sync.Mutex
}

// NewManagerRepository creates new repository for managers
func NewManagerRepository() (domain.ManagerRepository, error) {
	mr := managerRepository{
		managers: make(map[int]domain.Manager),
	}

	return &mr, nil
}

//
func (mr *managerRepository) Add(manager domain.Manager) error {
	if _, ok := mr.managers[manager.FplID]; ok {
		return errManagerAlreadyExists{
			fplID: manager.FplID,
		}
	}

	mr.Lock()
	mr.managers[manager.FplID] = manager
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
func (mr *managerRepository) GetByFplID(id int) (domain.Manager, error) {
	if manager, ok := mr.managers[id]; ok {
		return manager, nil
	}

	return domain.Manager{}, ErrManagerNotFound
}
