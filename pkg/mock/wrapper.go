package mock

import "fpl-live-tracker/pkg/wrapper"

type Wrapper struct {
	GetClubsFn        func() ([]wrapper.Club, error)
	GetFixturesFn     func() ([]wrapper.Fixture, error)
	GetGameweeksFn    func() ([]wrapper.Gameweek, error)
	GetPlayersFn      func() ([]wrapper.Player, error)
	GetPlayersStatsFn func(int) ([]wrapper.PlayerStats, error)
}

func (w *Wrapper) GetClubs() ([]wrapper.Club, error) {
	return w.GetClubsFn()
}

func (w *Wrapper) GetFixtures() ([]wrapper.Fixture, error) {
	return w.GetFixturesFn()
}

func (w *Wrapper) GetGameweeks() ([]wrapper.Gameweek, error) {
	return w.GetGameweeksFn()
}

func (w *Wrapper) GetPlayers() ([]wrapper.Player, error) {
	return w.GetPlayersFn()
}

func (w *Wrapper) GetPlayersStats(gameweekID int) ([]wrapper.PlayerStats, error) {
	return w.GetPlayersStatsFn(gameweekID)
}

func GetClubsOK() ([]wrapper.Club, error) {
	clubs := []wrapper.Club{
		{ID: 1, Name: "Arsenal", Shortname: "ARS"},
		{ID: 2, Name: "Aston Villa", Shortname: "AVL"},
		{ID: 3, Name: "Brentford", Shortname: "BRE"},
		{ID: 4, Name: "Brighton", Shortname: "BHA"},
		{ID: 5, Name: "Burnley", Shortname: "BUR"},
		{ID: 6, Name: "Chelsea", Shortname: "CHE"},
		{ID: 7, Name: "Crystal Palace", Shortname: "CRY"},
		{ID: 8, Name: "Everton", Shortname: "EVE"},
		{ID: 9, Name: "Leicester", Shortname: "LEI"},
		{ID: 10, Name: "Leeds", Shortname: "LEE"},
		{ID: 11, Name: "Liverpool", Shortname: "LIV"},
		{ID: 12, Name: "Man City", Shortname: "MCI"},
		{ID: 13, Name: "Man Utd", Shortname: "MUN"},
		{ID: 14, Name: "Newcastle", Shortname: "NEW"},
		{ID: 15, Name: "Norwich", Shortname: "NOR"},
		{ID: 16, Name: "Southampton", Shortname: "SOU"},
		{ID: 17, Name: "Spurs", Shortname: "TOT"},
		{ID: 18, Name: "Watford", Shortname: "WAT"},
		{ID: 19, Name: "West Ham", Shortname: "WHU"},
		{ID: 20, Name: "Wolves", Shortname: "WOL"},
	}

	return clubs, nil
}
