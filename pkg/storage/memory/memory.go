package memory

import tracker "fpl-live-tracker/pkg"

// type Manager struct {
// 	FplID    int
// 	FullName string
// }

type ManagerRepository struct {
	managers map[int]tracker.Manager
}

func (mr *ManagerRepository) Add(tracker.Manager) error {
	return nil
}

func (mr *ManagerRepository) AddMany(managers []tracker.Manager) error {
	return nil
}

func (mr *ManagerRepository) GetByFplID(id int) (*tracker.Manager, error) {
	return nil, nil
}
