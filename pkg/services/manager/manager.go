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
	Update() error
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

	err := ms.Update()
	if err != nil {
		log.Println("manager service: failed to init data", err)
		return nil, err
	}

	return &ms, nil
}

//
func (ms *managerService) Update() error {
	wrapperManager, err := ms.wr.GetManager(myID)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	manager := ms.convertToDomainManager(wrapperManager)
	err = ms.mr.Update(manager)
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

	for i := 0; i < len(manager.Team.Picks); i++ {
		tp := manager.Team.Picks[i]
		p, err := ms.ps.GetByID(tp.ID)
		if err != nil {
			log.Println("manager service:", err)
			continue
		}
		tp.Stats = p.Stats
		manager.Team.Picks[i] = tp
	}

	var totalPoints int
	for i := 0; i < 11; i++ {
		if manager.Team.Picks[i].IsCaptain {
			totalPoints += manager.Team.Picks[i].Stats.TotalPoints * 2
		} else {
			totalPoints += manager.Team.Picks[i].Stats.TotalPoints
		}
	}
	manager.Team.TotalPoints = totalPoints
	ms.mr.UpdateTeam(manager.ID, manager.Team)

	log.Println(manager)

	return nil
}

//
func (ms *managerService) GetByID(id int) (domain.Manager, error) {
	// TODO add validation
	return ms.mr.GetByID(id)
}

func (ms *managerService) convertToDomainManager(wm wrapper.Manager) domain.Manager {
	return domain.Manager{
		ID:       wm.ID,
		Name:     fmt.Sprintf("%s %s", wm.FirstName, wm.LastName),
		TeamName: wm.Name,
	}
}

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
