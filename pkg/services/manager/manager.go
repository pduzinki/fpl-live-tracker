package manager

import (
	"fmt"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"log"
)

// TODO remove later, and add support for handling more than one manager
var myID = 1239

type ManagerService interface {
	UpdateInfos() error
	UpdateTeams() error
	UpdatePoints() error
	GetByID(id int) (domain.Manager, error)
}

type managerService struct {
	mr domain.ManagerRepository
	ps player.PlayerService
	gs gameweek.GameweekService
	wr wrapper.Wrapper
}

//
func NewManagerService(mr domain.ManagerRepository, ps player.PlayerService,
	gs gameweek.GameweekService, wr wrapper.Wrapper) (ManagerService, error) {
	ms := managerService{
		mr: mr,
		ps: ps,
		gs: gs,
		wr: wr,
	}

	err := ms.UpdateInfos()
	if err != nil {
		log.Println("manager service: failed to init data", err)
		return nil, err
	}

	return &ms, nil
}

//
func (ms *managerService) UpdateInfos() error {
	wrapperManager, err := ms.wr.GetManager(myID)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	manager := ms.convertToDomainManager(wrapperManager)
	err = ms.mr.UpdateInfo(manager.ID, manager.Info)
	if err == storage.ErrManagerNotFound {
		err = ms.mr.Add(manager)
		if err != nil {
			log.Println("manager service:", err)
			return err
		}
	} else if err != nil {
		log.Println("manager service:", err)
		return err
	}

	return nil
}

//
func (ms *managerService) UpdateTeams() error {
	gameweek, err := ms.gs.GetCurrentGameweek()
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	wrapperTeam, err := ms.wr.GetManagersTeam(myID, gameweek.ID)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	team, err := ms.convertToDomainTeam(wrapperTeam)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	err = ms.mr.UpdateTeam(myID, team)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	return nil
}

//
func (ms *managerService) UpdatePoints() error {
	manager, err := ms.mr.GetByID(myID)
	if err != nil {
		return err
	}
	team := manager.Team

	err = ms.updateTeamStats(&team)
	if err != nil {
		return err
	}

	calculateTotalPoints(&team)
	calculateTotalPointsAfterSubs(&team)

	err = ms.mr.UpdateTeam(manager.ID, team)
	if err != nil {
		return err
	}

	log.Println(team.TotalPoints)
	log.Println(team.TotalPointsAfterSubs)

	return nil
}

//
func (ms *managerService) GetByID(id int) (domain.Manager, error) {
	manager := domain.Manager{ID: id}

	err := runManagerValidations(&manager, idHigherThanZero)
	if err != nil {
		return domain.Manager{}, err
	}

	return ms.mr.GetByID(id)
}

//
func (ms *managerService) convertToDomainManager(wm wrapper.Manager) domain.Manager {
	return domain.Manager{
		ID: wm.ID,
		Info: domain.ManagerInfo{
			Name:     fmt.Sprintf("%s %s", wm.FirstName, wm.LastName),
			TeamName: wm.Name,
		},
	}
}

//
func (ms *managerService) convertToDomainTeam(wt wrapper.Team) (domain.Team, error) {
	team := domain.Team{
		Picks: make([]domain.TeamPlayer, 0, 15),
	}

	for _, pick := range wt.Picks {
		p, err := ms.ps.GetByID(pick.ID)
		if err != nil {
			log.Println("manager service:", err)
			return domain.Team{}, err
		}

		dp := domain.TeamPlayer{
			Player:        p,
			IsCaptain:     pick.IsCaptain,
			IsViceCaptain: pick.IsViceCaptain,
		}

		team.Picks = append(team.Picks, dp)
	}

	return team, nil
}

//
func (ms *managerService) updateTeamStats(team *domain.Team) error {
	for i := 0; i < len(team.Picks); i++ {
		tp := team.Picks[i]
		p, err := ms.ps.GetByID(tp.ID)
		if err != nil {
			log.Println("manager service: failed to update manager's team stats", err)
			return err
		}
		tp.Stats = p.Stats
		team.Picks[i] = tp
	}

	return nil
}

//
func calculateTotalPoints(team *domain.Team) {
	var totalPoints int
	for i := 0; i < 11; i++ {
		if team.Picks[i].IsCaptain {
			totalPoints += team.Picks[i].Stats.TotalPoints * 2
		} else {
			totalPoints += team.Picks[i].Stats.TotalPoints
		}
	}
	team.TotalPoints = totalPoints
}

//
func calculateTotalPointsAfterSubs(team *domain.Team) {
	totalPointsAfterSubs := team.TotalPoints

	// get formation
	// formation := getFormation(team)
	// log.Println(formation)

	// get formation with taking into account the missing players
	realFormation, players := getSubFormationAndSubPlayers(team)
	// log.Println(realFormation, players)

	subs := getSubs(team, realFormation, players)

	for _, s := range subs {
		totalPointsAfterSubs += s.Stats.TotalPoints
	}

	team.TotalPointsAfterSubs = totalPointsAfterSubs
}

//
func getSubFormationAndSubPlayers(team *domain.Team) ([4]int, []domain.TeamPlayer) {
	var gkps, defs, mids, fwds int
	toBeSubbed := make([]domain.TeamPlayer, 0)

	for i := 0; i < 11; i++ {
		player := &team.Picks[i]

		if needsSubstitution(player) {
			toBeSubbed = append(toBeSubbed, *player)
		} else {
			if player.Info.Position == "GKP" {
				gkps++
			} else if player.Info.Position == "DEF" {
				defs++
			} else if player.Info.Position == "MID" {
				mids++
			} else if player.Info.Position == "FWD" {
				fwds++
			}
		}
	}

	return [4]int{gkps, defs, mids, fwds}, toBeSubbed
}

//
func needsSubstitution(player *domain.TeamPlayer) bool {
	stats := player.Stats

	var atLeastOneFixtureStarted bool
	for _, f := range stats.FixturesInfo {
		if f.Started {
			atLeastOneFixtureStarted = true
			break
		}
	}

	if stats.Minutes == 0 && atLeastOneFixtureStarted {
		return true
	}
	return false
}

func getSubs(team *domain.Team, formation [4]int, players []domain.TeamPlayer) []domain.TeamPlayer {
	subs := make([]domain.TeamPlayer, 0)

	for _, p := range players {
		_ = p
		if formation[0] == 0 {
			// goalkeeper sub
			subs = append(subs, team.Picks[11])
		}

		if formation[1] < 3 {
			// need a defender first
			for _, pp := range team.Picks[11:] {
				if pp.Info.Position == "DEF" {
					subs = append(subs, pp)
					formation[1]++
				}
			}
		}

		if formation[3] < 1 {
			// need a forward first
			for _, pp := range team.Picks[11:] {
				if pp.Info.Position == "FWD" {
					subs = append(subs, pp)
					formation[1]++
				}
			}
		}
	}
	return subs
}
