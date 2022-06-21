package domain

//
type Team struct {
	ID         int    `bson:"_id"`
	GameweekID int    `bson:"GameweekID"`
	ActiveChip string `bson:"ActiveChip"`
	// TODO remove redundant GwPointsXyz fields later
	GwPoints                int          `bson:"GwPoints"`                // points gained in current gw
	GwPointsWithSubs        int          `bson:"GwPointsWithSubs"`        // points gained in current gw with possible subs included
	GwPointsWithHits        int          `bson:"GwPointsWithHits"`        // points gained in current gw with hits subtracted
	GwPointsWithHitsAndSubs int          `bson:"GwPointsWithHitsAndSubs"` // points gained in current gw with hits subtracted, with possible subs included
	GwHitPoints             int          `bson:"GwHitPoints"`             // hit points taken in current gw
	PrevOverallPoints       int          `bson:"PrevOverallPoints"`       // overall points before start of current gw
	OverallPoints           int          `bson:"OverallPoints"`           // overall points with current gw points included
	GwRank                  int          `bson:"GwRank"`
	OverallRank             int          `bson:"OverallRank"`
	Picks                   []TeamPlayer `bson:"Picks"`
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
