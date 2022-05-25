package domain

// Manager represents a human being that plays Fantasy Premier League
type Manager struct {
	ID   int         `bson:"_id"`
	Info ManagerInfo `bson:"ManagerInfo"`
	Team Team        `bson:"Team"` // TODO remove Team from Manager
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
	UpdateInfo(managerID int, info ManagerInfo) error
	UpdateTeam(managerID int, team Team) error // TODO remove this
	GetByID(id int) (Manager, error)
	GetCount() (int, error)
}
