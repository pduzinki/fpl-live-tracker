package mock

import "fpl-live-tracker/pkg/wrapper"

type Wrapper struct {
	GetClubsFn         func() ([]wrapper.Club, error)
	GetFixturesFn      func() ([]wrapper.Fixture, error)
	GetGameweeksFn     func() ([]wrapper.Gameweek, error)
	GetPlayersFn       func() ([]wrapper.Player, error)
	GetPlayersStatsFn  func(int) ([]wrapper.PlayerStats, error)
	GetManagersCountFn func() (int, error)
	GetManagerFn       func(id int) (wrapper.Manager, error)
	GetTeamFn          func(managerID, gameweekID int) (wrapper.Team, error)
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

func (w *Wrapper) GetManagersCount() (int, error) {
	return w.GetManagersCountFn()
}

func (w *Wrapper) GetManager(id int) (wrapper.Manager, error) {
	return w.GetManagerFn(id)
}
func (w *Wrapper) GetTeam(managerID, gameweekID int) (wrapper.Team, error) {
	return w.GetTeamFn(managerID, gameweekID)
}
