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
func (mr *managerRepository) Add(manager *tracker.Manager) error {
	return nil
}

//
func (mr *managerRepository) AddMany(managers []tracker.Manager) error {
	return nil
}

//
func (mr *managerRepository) GetByFplID(id int) (*tracker.Manager, error) {
	return nil, nil
}
