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

const tripleCaptainActive = "3xc"
const benchBoostActive = "bboost"

// TODO remove later, and add support for handling more than a few managers
// var IDs []int = []int{1239, 445331, 1056968, 2037831}
var myID = 1239

type ManagerService interface {
	AddNew() error
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
func (ms *managerService) AddNew() error {
	fplManagers, err := ms.wr.GetManagersCount()
	if err != nil {
		return err
	}

	storageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	if fplManagers > storageManagers {
		log.Println("manager service: adding new managers...")
		// TODO
	}

	return fmt.Errorf("not implemented")
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

	manager, err := ms.GetByID(myID)
	if err != nil {
		log.Println("manager service:", err)
	}

	if manager.Team.GameweekID == gameweek.ID {
		// log.Println("manager service: team already up-to-date")
		return nil
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

	err = ms.updateTeamPlayersStats(&team)
	if err != nil {
		return err
	}

	totalPoints := calculateTotalPoints(&team)
	subPoints := calculateSubPoints(&team)

	team.TotalPoints = totalPoints - team.HitPoints
	team.TotalPointsAfterSubs = totalPoints + subPoints - team.HitPoints

	// log.Println(team.TotalPoints)
	// log.Println(team.TotalPointsAfterSubs)

	err = ms.mr.UpdateTeam(manager.ID, team)
	if err != nil {
		return err
	}

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
		GameweekID: wt.EntryHistory.GameweekID,
		Picks:      make([]domain.TeamPlayer, 0, 15),
		ActiveChip: wt.ActiveChip,
		HitPoints:  wt.EntryHistory.EventTransfersCost,
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
func (ms *managerService) updateTeamPlayersStats(team *domain.Team) error {
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
