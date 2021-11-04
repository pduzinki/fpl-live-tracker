package tracker

// Manager represents a human being that plays Fantasy Premier League
type Manager struct {
	FplID    int
	FullName string
	TeamName string
	// TODO add more fields
}

// ManagerRepository is an interface for interacting with Manager storage
type ManagerRepository interface {
	Add(manager Manager) error
	AddMany(managers []Manager) error
	GetByFplID(id int) (Manager, error)
}
