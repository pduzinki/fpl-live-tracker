package domain

// Manager represents a human being that plays Fantasy Premier League
type Manager struct {
	ID       int
	Name     string
	TeamName string
	Team     Team
}

// ManagerRepository is an interface for interacting with Manager storage
type ManagerRepository interface {
	Add(manager Manager) error
	AddMany(managers []Manager) error
	GetByFplID(id int) (Manager, error)
}

//
type Team struct {
	Picks []TeamPlayer
	/*
		gk:
		defs:
		mids:
		fwds:
		benchGk:
		bench:
	*/
	TotalPoints int
	// OverallRank int
}

//
type TeamPlayer struct {
	Player
	IsCaptain     bool
	IsViceCaptain bool
}
