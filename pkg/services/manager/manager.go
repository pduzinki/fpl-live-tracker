package manager

import (
	"fmt"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/wrapper"
	"log"
)

type ManagerService interface {
	Update() error
	// UpdateTeams() error
	// UpdatePoints() error
	GetByID(id int) (domain.Manager, error)
}

type managerService struct {
	mr domain.ManagerRepository
	ps player.PlayerService
	wr wrapper.Wrapper
}

//
func NewManagerService(mr domain.ManagerRepository, ps player.PlayerService,
	wr wrapper.Wrapper) (ManagerService, error) {
	ms := managerService{
		mr: mr,
		ps: ps,
		wr: wr,
	}

	err := ms.Update()
	if err != nil {
		log.Println("manager service: failed to init data", err)
		return nil, err
	}

	return &ms, nil
}

func (ms *managerService) Update() error {
	// TODO all this is temporary
	myID := 1239

	wrapperManager, err := ms.wr.GetManager(myID)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	wrapperTeam, err := ms.wr.GetManagersTeam(myID, 23)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	manager := ms.convertToDomainManager(wrapperManager)
	team, err := ms.convertToDomainTeam(wrapperTeam)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	log.Println(manager)
	for _, p := range team.Picks {
		log.Println(p)
	}

	return nil
}

//
func (ms *managerService) GetByID(id int) (domain.Manager, error) {
	return domain.Manager{}, nil
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
