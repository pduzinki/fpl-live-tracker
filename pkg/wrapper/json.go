package wrapper

type Manager struct {
	ID        int    `json:"id"`
	FirstName string `json:"player_first_name"`
	LastName  string `json:"player_last_name"`
	Name      string `json:"name"`
}

type Team struct {
	// TODO
}

type Bootstrap struct {
	// Count int        `json:"total_players"`
	Gameweeks []Gameweek `json:"events"`
	Clubs     []Club     `json:"teams"`
	Players   []Player   `json:"elements"`
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
	Event int `json:"event"`
	ID    int `json:"id"`
	// Started             bool `json:"started"`
	// FinishedProvisional bool `json:"finished_provisional"`
	// Finished            bool `json:"finished"`
	TeamA int `json:"team_a"`
	TeamH int `json:"team_h"`
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

// TODO rename structs below to something more meaningful
type Elements struct {
	Elements []PlayerLive `json:"elements"`
}

type PlayerLive struct {
	ID    int         `json:"id"`
	Stats PlayerStats `json:"stats"`
}

type PlayerStats struct {
	Minutes     int `json:"minutes"`
	TotalPoints int `json:"total_point"`
}
