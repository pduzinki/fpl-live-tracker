package wrapper

type Manager struct {
	ID        int    `json:"id"`
	FirstName string `json:"player_first_name"`
	LastName  string `json:"player_last_name"`
	Name      string `json:"name"`
}

type Team struct {
	ActiveChip   string       `json:"active_chip"`
	EntryHistory EntryHistory `json:"entry_history"`
	Picks        []Pick       `json:"picks"`
}

type EntryHistory struct {
	GameweekID         int `json:"event"`
	Points             int `json:"points"`
	TotalPoints        int `json:"total_points"`
	OverallRank        int `json:"overall_rank"`
	EventTransfers     int `json:"event_transfers"`
	EventTransfersCost int `json:"event_transfers_cost"`
}

type Pick struct {
	ID            int  `json:"element"`
	Position      int  `json:"position"`
	Multiplier    int  `json:"multiplier"`
	IsCaptain     bool `json:"is_captain"`
	IsViceCaptain bool `json:"is_vice_captain"`
}

type Bootstrap struct {
	ManagersCount int        `json:"total_players"`
	Gameweeks     []Gameweek `json:"events"`
	Clubs         []Club     `json:"teams"`
	Players       []Player   `json:"elements"`
}

type Gameweek struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Finished     bool   `json:"finished"`
	IsCurrent    bool   `json:"is_current"`
	IsNext       bool   `json:"is_next"`
	DeadlineTime string `json:"deadline_time"`
}

type Fixture struct {
	Event               int           `json:"event"`
	ID                  int           `json:"id"`
	KickoffTime         string        `json:"kickoff_time"`
	Started             bool          `json:"started"`
	Finished            bool          `json:"finished"`
	FinishedProvisional bool          `json:"finished_provisional"`
	TeamA               int           `json:"team_a"`
	TeamH               int           `json:"team_h"`
	Stats               []FixtureStat `json:"stats"`
}

type FixtureStat struct {
	Identifier string             `json:"identifier"`
	TeamA      []FixtureStatValue `json:"a"`
	TeamH      []FixtureStatValue `json:"h"`
}

type FixtureStatValue struct {
	Value   int `json:"value"`
	Element int `json:"element"`
}

type Club struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Shortname string `json:"short_name"`
}

type Player struct {
	ID       int    `json:"id"`
	Team     int    `json:"team"`
	Position int    `json:"element_type"`
	WebName  string `json:"web_name"`
}

type Elements struct {
	PlayersStats []PlayerStats `json:"elements"`
}

type PlayerStats struct {
	ID      int       `json:"id"`
	Stats   Stats     `json:"stats"`
	Explain []Explain `json:"explain"`
}

type Stats struct {
	Minutes     int `json:"minutes"`
	TotalPoints int `json:"total_points"`
}

type Explain struct {
	Fixture int `json:"fixture"`
	// more fields not needed for
}
