package memory

import (
	"errors"
	tracker "fpl-live-tracker/pkg"
	"sync"
)

// type Manager struct {
// 	FplID    int
// 	FullName string
// }

var ErrRecordAlreadyExists error = errors.New("storage: record already exists") // TODO not sure if that's the best place for those
var ErrRecordNotFound error = errors.New("storage: record not found")

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
		return ErrRecordAlreadyExists
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
			return err // TODO return a more appropriate error (perhaps add custom type error)
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
