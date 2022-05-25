package domain

//
type Team struct {
	ID                   int          `bson:"_id"`
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

// TeamRepository is an interface for interacting with Team storage
type TeamRepository interface {
	Add(Team Team) error
	Update(ID int, team Team) error // TODO probably remove ID arg
	GetByID(id int) (Team, error)
	GetCount() (int, error)
}
