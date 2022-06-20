package domain

//
type Team struct {
	ID              int          `bson:"_id"`
	GameweekID      int          `bson:"GameweekID"`
	Picks           []TeamPlayer `bson:"Picks"`
	ActiveChip      string       `bson:"ActiveChip"`
	HitPoints       int          `bson:"HitPoints"`
	Points          int          `bson:"TotalPoints"`
	PointsAfterSubs int          `bson:"TotalPointsAfterSubs"`
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

// TeamRepository is an interface for interacting with Team storage
type TeamRepository interface {
	Add(team Team) error
	Update(team Team) error
	GetByID(ID int) (Team, error)
	GetCount() (int, error)
}
