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
		Picks:      make([]domain.TeamPlayer, 0, 15),
		ActiveChip: wt.ActiveChip,
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
	captainMultiplier := 2
	playersCount := 11

	if team.ActiveChip == "3xc" { // triple captain
		captainMultiplier = 3
	} else if team.ActiveChip == "bboost" { // bench boost
		playersCount = 15
	}

	var didCaptainPlay bool

	var totalPoints int
	for i := 0; i < playersCount; i++ {
		if team.Picks[i].IsCaptain {
			totalPoints += team.Picks[i].Stats.TotalPoints * captainMultiplier
			didCaptainPlay = played(&team.Picks[i])
		} else {
			totalPoints += team.Picks[i].Stats.TotalPoints
		}
	}

	if !didCaptainPlay {
		for i := 0; i < playersCount; i++ {
			if team.Picks[i].IsViceCaptain {
				totalPoints += team.Picks[i].Stats.TotalPoints * (captainMultiplier - 1)
			}
		}
	}

	team.TotalPoints = totalPoints
}

//
func calculateTotalPointsAfterSubs(team *domain.Team) {
	/*
		(legit formation == 1 gkp, at least 3 defs, and at least 1 fwd)
		get live formation

		for p in range bench
			if too few defs, sub only if p is def
			if too few fwds, sub only if p is fwd
			else sub p
	*/
	if team.ActiveChip == "bboost" {
		return
	}

	totalPointsAfterSubs := team.TotalPoints

	liveFormation := getLiveFormation(team)
	bench := team.Picks[12:]

	// check if goalkeeper needs a substitution
	if liveFormation[0] == 0 {
		benchGk := &team.Picks[11]
		if played(benchGk) {
			totalPointsAfterSubs += benchGk.Stats.TotalPoints
			log.Println("IN:", benchGk.Info.Name)
		}
	}

	// check if outfield players need substitutions
	subsNeeded := 10 - (liveFormation[1] + liveFormation[2] + liveFormation[3])
	subsIn := make([]domain.TeamPlayer, 0)

	for _, b := range bench {
		if subsNeeded == 0 {
			break
		}
		pos := b.Info.Position

		if liveFormation[1] < 3 { // too few defs, add only if b is def
			if pos == "DEF" && played(&b) {
				subsIn = append(subsIn, b)
				liveFormation[1]++
				subsNeeded--
			}
			continue
		}

		if liveFormation[3] < 1 { // too few fwds, add only if b is fwd
			if pos == "FWD" && played(&b) {
				subsIn = append(subsIn, b)
				liveFormation[3]++
				subsNeeded--
			}
			continue
		}

		if played(&b) {
			subsNeeded--
			subsIn = append(subsIn, b)
		}
	}

	for _, s := range subsIn {
		log.Println("IN:", s.Info.Name)
		totalPointsAfterSubs += s.Stats.TotalPoints
	}

	team.TotalPointsAfterSubs = totalPointsAfterSubs
}

//
func getLiveFormation(team *domain.Team) [4]int {
	var gkps, defs, mids, fwds int

	for _, p := range team.Picks[:11] {
		pos := p.Info.Position
		if pos == "GKP" && played(&p) {
			gkps++
		} else if pos == "DEF" && played(&p) {
			defs++
		} else if pos == "MID" && played(&p) {
			mids++
		} else if pos == "FWD" && played(&p) {
			fwds++
		}
	}

	return [4]int{gkps, defs, mids, fwds}
}

// played returns boolean value, indicating if given player
// had played some minutes in his fixtures
func played(player *domain.TeamPlayer) bool {
	stats := player.Stats

	var fixtureStarted bool
	for _, fixture := range stats.FixturesInfo {
		if fixture.Started {
			fixtureStarted = true
			break
		}
	}

	if stats.Minutes == 0 && fixtureStarted {
		return false
	}
	return true
}
