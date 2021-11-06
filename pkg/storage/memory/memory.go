package memory

import (
	tracker "fpl-live-tracker/pkg"
	"sync"
)

// type Manager struct {
// 	FplID    int
// 	FullName string
// }

//
type managerRepository struct {
	managers map[int]tracker.Manager
	sync.Mutex
}

// NewManagerRepository creates new repository for managers
func NewManagerRepository() (tracker.ManagerRepository, error) {
	mr := managerRepository{
		managers: make(map[int]tracker.Manager),
	}

	return &mr, nil
}

//
func (mr *managerRepository) Add(manager tracker.Manager) error {
	if _, ok := mr.managers[manager.FplID]; ok {
		return errorRecordAlreadyExists{
			fplID: manager.FplID,
		}
	}

	mr.Lock()
	mr.managers[manager.FplID] = manager
	mr.Unlock()

	return nil
}

//
func (mr *managerRepository) AddMany(managers []tracker.Manager) error {
	for _, manager := range managers {
		err := mr.Add(manager)
		if err != nil {
			return err
		}
	}
	return nil
}

//
func (mr *managerRepository) GetByFplID(id int) (tracker.Manager, error) {
	if manager, ok := mr.managers[id]; ok {
		return manager, nil
	}

	return tracker.Manager{}, ErrRecordNotFound
}
