package domain

// Manager represents a human being that plays Fantasy Premier League
type Manager struct {
	ID          int         `bson:"_id"`
	UpdatedInGw int         `bson:"UpdatedInGw"`
	Info        ManagerInfo `bson:"ManagerInfo"`
}

//
type ManagerInfo struct {
	Name     string `bson:"Name"`
	TeamName string `bson:"TeamName"` // TODO move fields from ManagerInfo to Manager
}

// ManagerRepository is an interface for interacting with Manager storage
type ManagerRepository interface {
	Add(manager Manager) error
	AddMany(managers []Manager) error
	Update(manager Manager) error
	GetByID(id int) (Manager, error)
	GetCount() (int, error)
}
