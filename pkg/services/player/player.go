package player

import (
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/wrapper"
	"log"

	domain "fpl-live-tracker/pkg"
)

var positions = map[int]string{
	1: "GKP",
	2: "DEF",
	3: "MID",
	4: "FWD",
}

type PlayerService interface {
	Update() error
	UpdateStats() error
	GetByID(ID int) (domain.Player, error)
}

type playerService struct {
	wrapper wrapper.Wrapper
	pr      domain.PlayerRepository
	cs      club.ClubService
	gs      gameweek.GameweekService
}

//
func NewPlayerService(w wrapper.Wrapper, pr domain.PlayerRepository, cs club.ClubService, gs gameweek.GameweekService) PlayerService {
	return &playerService{
		wrapper: w,
		pr:      pr,
		cs:      cs,
		gs:      gs,
	}
}

//
func (ps *playerService) Update() error {
	wrapperPlayers, err := ps.wrapper.GetPlayers()
	if err != nil {
		return err
	}

	players := make([]domain.Player, len(wrapperPlayers))
	for i, wp := range wrapperPlayers {
		club, err := ps.cs.GetClubByID(wp.Team)
		if err != nil {
			log.Println(err)
		}

		players[i] = domain.Player{
			ID:       wp.ID,
			Name:     wp.WebName,
			Position: positions[wp.Position],
			Club:     club,
		}

		err = ps.pr.Add(players[i])
		if err != nil {
			log.Println(err)
		}
	}

	return nil
}

//
func (ps *playerService) UpdateStats() error {
	gw, err := ps.gs.GetCurrentGameweek()
	if err != nil {
		return err
	}

	liveData, err := ps.wrapper.GetPlayersStats(gw.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	// log.Println(liveData[13])

	for _, item := range liveData {
		_ = item
	}

	return nil
}

//
func (ps *playerService) GetByID(ID int) (domain.Player, error) {
	// TODO
	return domain.Player{}, nil
}
