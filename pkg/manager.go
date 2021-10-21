package tracker

//
type Manager struct {
	FplID    int
	FullName string
}

//
type ManagerRepository interface {
	Add(manager *Manager) error
	AddMany(managers []Manager) error
	GetByFplID(id int) (*Manager, error)
}
