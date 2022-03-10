package domain

// Manager represents a human being that plays Fantasy Premier League
type Manager struct {
	ID   int         `bson:"_id"`
	Info ManagerInfo `bson:"ManagerInfo"`
	Team Team        `bson:"Team"`
}

//
type ManagerInfo struct {
	Name     string `bson:"Name"`
	TeamName string `bson:"TeamName"`
}

//
type Team struct {
	GameweekID           int          `bson:"GameweekID"`
	Picks                []TeamPlayer `bson:"Picks"`
	ActiveChip           string       `bson:"ActiveChip"`
	HitPoints            int          `bson:"HitPoints"`
	TotalPoints          int          `bson:"TotalPoints"`
	TotalPointsAfterSubs int          `bson:"TotalPointsAfterSubs"`
	// OverallRank int
}

//
type TeamPlayer struct {
	Player        `bson:"Player"`
	IsCaptain     bool `bson:"IsCaptain"`
	IsViceCaptain bool `bson:"IsViceCaptain"`
	SubIn         bool `bson:"SubIn"`
	// SubOut        bool
}

// ManagerRepository is an interface for interacting with Manager storage
type ManagerRepository interface {
	Add(manager Manager) error
	AddMany(managers []Manager) error
	UpdateInfo(managerID int, info ManagerInfo) error
	UpdateTeam(managerID int, team Team) error
	GetByID(id int) (Manager, error)
	GetCount() (int, error)
}
