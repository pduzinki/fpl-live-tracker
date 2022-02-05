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

// PlayerService is an interface for interacting with players
type PlayerService interface {
	Update() error
	UpdateStats() error
	GetByID(ID int) (domain.Player, error)
	GetAll() ([]domain.Player, error)
}

// playerService implements PlayerService interface
type playerService struct {
	wrapper wrapper.Wrapper
	pr      domain.PlayerRepository
	cs      club.ClubService
	fs      fixture.FixtureService
	gs      gameweek.GameweekService
}

// bonusPlayer is a helper struct used in process of calculating predicted bonus points
type bonusPlayer struct {
	playerID    int
	bonusPoints int
}

// NewPlayerService new instance of PlayerService, and fills
// underlying data storage with data from FPL API
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

// Update queries FPL API and updates all players basic data
// (i.e. name, position, club), in its underlying player storage
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

// UpdateStats queries FPL API and updates all players current gameweek
// stats data in its underlying player storage
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

	err = ps.updatePredictedBonusPoints()
	if err != nil {
		log.Println("player service: failed to update predicted bonus points", err)
	}

	return nil
}

// GetByID returns player with given ID, or error otherwise
func (ps *playerService) GetByID(ID int) (domain.Player, error) {
	player := domain.Player{ID: ID}

	err := runPlayerValidations(&player, idHigherThanZero)
	if err != nil {
		return domain.Player{}, err
	}

	return ps.pr.GetByID(ID)
}

// GetAll returns slice of all existing players
func (ps *playerService) GetAll() ([]domain.Player, error) {
	return ps.pr.GetAll()
}

// convertToDomainPlayer returns domain.Player, consistent with given wrapper.Player
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

// convertToDomainPlayerStats returns domain.PlayerStats,
// consistent with given wrapper.PlayerStats
func (ps *playerService) convertToDomainPlayerStats(ws wrapper.PlayerStats) domain.PlayerStats {
	return domain.PlayerStats{
		Minutes:     ws.Stats.Minutes,
		TotalPoints: ws.Stats.TotalPoints,
	}
}

// updatePredictedBonusPoints add predicted bonus points for players,
// which bonus points weren't confirmed yet
func (ps *playerService) updatePredictedBonusPoints() error {
	gw, err := ps.gs.GetCurrentGameweek()
	if err != nil {
		return err
	}

	liveFixtures, err := ps.fs.GetLiveFixtures(gw.ID)
	if err != nil {
		log.Println("player service:", err)
		return err
	}

	for _, f := range liveFixtures {
		bpsStats, ok := f.Stats["bps"]
		if !ok {
			log.Println("player service: bps stats not found in live fixture")
			continue
		}

		// merge stats of home and away players together, and sort them in descending order
		allPlayersStats := make([]domain.FixtureStatPair, len(bpsStats.AwayPlayersStats)+len(bpsStats.HomePlayersStats))
		allPlayersStats = append(bpsStats.AwayPlayersStats, bpsStats.HomePlayersStats...)
		sort.Slice(allPlayersStats, func(i, j int) bool {
			return (allPlayersStats[i].Value > allPlayersStats[j].Value)
		})

		topBPS := findTopBPS(allPlayersStats)
		bp := findPlayersAndBonusPoints(allPlayersStats, topBPS)
		for _, pair := range bp {
			err := ps.addBonusPoints(pair.playerID, pair.bonusPoints)
			if err != nil {
				log.Println("player service: failed to add bonus points", err)
			}
		}
	}

	return nil
}

// findTopBPS returns slice of current top 3 bps values in playerStats
func findTopBPS(playersStats []domain.FixtureStatPair) []int {
	bpsToReward := make([]int, 0, 3)
	bpsToReward = append(bpsToReward, playersStats[0].Value)
	for _, p := range playersStats {
		if len(bpsToReward) == 3 {
			break
		}
		if p.Value == bpsToReward[len(bpsToReward)-1] {
			continue
		}
		bpsToReward = append(bpsToReward, p.Value)
	}

	return bpsToReward
}

// findPlayersAndBonusPoints returns slice of bonusPlayer, which holds ID of player,
// and amount of predicted bonus points to be added
func findPlayersAndBonusPoints(playersStats []domain.FixtureStatPair, topBPS []int) []bonusPlayer {
	bp := make([]bonusPlayer, 0)

	bonus := 3
	awardedPlayersCounts := 0
	for i := 0; i < 3; i++ {
		for _, p := range playersStats {
			if p.Value == topBPS[i] {
				bp = append(bp, bonusPlayer{p.PlayerID, bonus})
				awardedPlayersCounts++
			}
		}
		if awardedPlayersCounts >= 3 {
			break
		}
		bonus--
	}

	return bp
}

// addBonusPoints adds predicted bonus points to player with given ID
func (ps *playerService) addBonusPoints(playerID, points int) error {
	player, err := ps.pr.GetByID(playerID)
	if err != nil {
		return err
	}
	player.Stats.TotalPoints += points
	err = ps.pr.UpdateStats(playerID, player.Stats)
	if err != nil {
		return err
	}

	return nil
}
