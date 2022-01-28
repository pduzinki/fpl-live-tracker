package player

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/services/fixture"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"sort"
)

type PlayerService interface {
	Update() error
	UpdateStats() error
	GetByID(ID int) (domain.Player, error)
}

type playerService struct {
	wrapper wrapper.Wrapper
	pr      domain.PlayerRepository
	cs      club.ClubService
	fs      fixture.FixtureService
	gs      gameweek.GameweekService
}

//
func NewPlayerService(w wrapper.Wrapper, pr domain.PlayerRepository, cs club.ClubService,
	fs fixture.FixtureService, gs gameweek.GameweekService) (PlayerService, error) {
	ps := playerService{
		wrapper: w,
		pr:      pr,
		cs:      cs,
		fs:      fs,
		gs:      gs,
	}

	err := ps.Update()
	if err != nil {
		return nil, err
	}

	return &ps, nil
}

//
func (ps *playerService) Update() error {
	wrapperPlayers, err := ps.wrapper.GetPlayers()
	if err != nil {
		return err
	}

	players := make([]domain.Player, len(wrapperPlayers))
	for i, wp := range wrapperPlayers {
		p, err := ps.convertToDomainPlayer(wp)
		if err != nil {
			log.Println("player service:", err)
			return err
		}
		players[i] = p
	}

	for _, p := range players {
		err := ps.pr.Update(p)
		if err == storage.ErrPlayerNotFound {
			err = ps.pr.Add(p)
			if err != nil {
				log.Println("player service:", err)
				return err
			}
		} else if err != nil {
			log.Println("player service:", err)
			return err
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

	playersStats, err := ps.wrapper.GetPlayersStats(gw.ID)
	if err != nil {
		log.Println(err)
		return err
	}

	for _, ws := range playersStats {
		stats := ps.convertToDomainPlayerStats(ws)
		err := ps.pr.UpdateStats(ws.ID, stats)
		if err != nil {
			log.Println("player service: failed to update player stats", err)
		}
	}

	// add live bonus points
	liveFixtures, err := ps.fs.GetLiveFixtures(gw.ID)
	if err != nil {
		log.Println("player service:", err)
	}

	for _, f := range liveFixtures {
		s, ok := f.Stats["bps"]
		if !ok {
			log.Println("player service: bps stats not found in live fixture")
			continue
		}

		allPlayersStats := make([]domain.FixtureStatValue, len(s.AwayPlayersStats)+len(s.HomePlayersStats))
		allPlayersStats = append(s.AwayPlayersStats, s.HomePlayersStats...)

		// sort in descending order
		sort.Slice(allPlayersStats, func(i, j int) bool {
			return (allPlayersStats[i].Value > allPlayersStats[j].Value)
		})

		topBPS := ps.topBPS(allPlayersStats)
		ps.awardBonusPoints(allPlayersStats, topBPS)
	}

	return nil
}

//
func (ps *playerService) GetByID(ID int) (domain.Player, error) {
	player := domain.Player{ID: ID}

	err := runPlayerValidations(&player, idHigherThanZero)
	if err != nil {
		return domain.Player{}, err
	}

	return ps.pr.GetByID(ID)
}

func (ps *playerService) topBPS(bpsPlayers []domain.FixtureStatValue) []int {
	topBPS := make([]int, 0, 3)
	topBPS = append(topBPS, bpsPlayers[0].Value)
	for _, p := range bpsPlayers {
		if len(topBPS) == 3 {
			break
		}
		if p.Value == topBPS[len(topBPS)-1] {
			continue
		}
		topBPS = append(topBPS, p.Value)
	}

	return topBPS
}

func (ps *playerService) awardBonusPoints(allPlayersStats []domain.FixtureStatValue, topBPS []int) {
	log.Println("----")
	bonus := 3
	awardedPlayersCounts := 0
	for i := 0; i < 3; i++ {
		for _, p := range allPlayersStats {
			if p.Value == topBPS[i] {
				ps.addBPS(p.PlayerID, bonus)
				log.Printf("playerID: %d, bps: %d, bonus %d", p.PlayerID, p.Value, bonus)
				awardedPlayersCounts++
			}
		}
		if awardedPlayersCounts >= 3 {
			break
		}
		bonus--
	}
}

func (ps *playerService) addBPS(playerID, points int) {
	player, _ := ps.pr.GetByID(playerID) // TODO check error
	player.Stats.TotalPoints += points
	log.Printf("playerID: %d name: %s bonus: %d", player.ID, player.Name, points)
	ps.pr.UpdateStats(playerID, player.Stats)
}

func (ps *playerService) convertToDomainPlayer(wp wrapper.Player) (domain.Player, error) {
	club, err := ps.cs.GetClubByID(wp.Team)
	if err != nil {
		log.Println(err)
		return domain.Player{}, err
	}

	return domain.Player{
		ID:       wp.ID,
		Name:     wp.WebName,
		Position: domain.PlayerPosition[wp.Position],
		Club:     club,
	}, nil
}

func (ps *playerService) convertToDomainPlayerStats(ws wrapper.PlayerStats) domain.PlayerStats {
	return domain.PlayerStats{
		Minutes:     ws.Stats.Minutes,
		TotalPoints: ws.Stats.TotalPoints,
	}
}
