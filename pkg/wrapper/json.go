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

type League struct {
	// TODO
}

type Bootstrap struct {
	// Count int        `json:"total_players"`
	Gws []Gameweek `json:"events"`
}

type Gameweek struct {
	ID           int    `json:"id"`
	Name         string `json:"name"`
	Finished     bool   `json:"finished"`
	IsCurrent    bool   `json:"is_current"`
	IsNext       bool   `json:"is_next"`
	DeadlineTime string `json:"deadline_time"`
}
