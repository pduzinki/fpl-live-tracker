package domain

// Manager represents a human being that plays Fantasy Premier League
type Manager struct {
	ID   int
	Info ManagerInfo
	Team Team
}

//
type ManagerInfo struct {
	Name     string
	TeamName string
}

//
type Team struct {
	Picks                []TeamPlayer
	ActiveChip           string
	HitPoints            int
	TotalPoints          int
	TotalPointsAfterSubs int
	// OverallRank int
}

//
type TeamPlayer struct {
	Player
	IsCaptain     bool
	IsViceCaptain bool
	SubIn         bool
	// SubOut        bool
}

// ManagerRepository is an interface for interacting with Manager storage
type ManagerRepository interface {
	Add(manager Manager) error
	AddMany(managers []Manager) error
	UpdateInfo(managerID int, info ManagerInfo) error
	UpdateTeam(managerID int, team Team) error
	GetByID(id int) (Manager, error)
}
