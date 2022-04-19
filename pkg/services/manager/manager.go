package manager

import (
	"fmt"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"runtime"
	"sync"
	"time"
)

const tripleCaptainActive = "3xc"
const benchBoostActive = "bboost"

// TODO remove later, and add support for handling more than a few managers
// var IDs []int = []int{1239, 445331, 1056968, 2037831}
// var myID = 1239

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

	return &ms, nil
}

//
func (ms *managerService) AddNew() error {
	log.Println("manager service: AddNew started")
	inFplManagers, err := ms.wr.GetManagersCount()
	if err != nil {
		return err
	}

	inStorageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	if inFplManagers == inStorageManagers {
		return nil // everything up-to-date, nothing to do here
	}

	chanSize := runtime.NumCPU() * 4
	workerCount := runtime.NumCPU() * 16

	ids := make(chan int, chanSize)
	failed := make(chan int, chanSize)
	managers := make(chan wrapper.Manager, chanSize)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	newManagersCount := inFplManagers - inStorageManagers
	workerWg.Add(newManagersCount)

	for i := 0; i <= workerCount; i++ {
		go func() {
			for id := range ids {
				wm, err := ms.wr.GetManager(id)
				if err != nil { // TODO improve error handling
					failed <- id
					time.Sleep(10 * time.Second)
				}

				managers <- wm
				workerWg.Done()
			}
		}()
	}

	innerWg.Add(1)
	go func() {
		// send to ids chan
		for id := inStorageManagers + 1; id <= inFplManagers; id++ {
			ids <- id
		}
		innerWg.Done()
	}()

	innerWg.Add(1)
	go func() {
		// receive from failed chan, send to ids chan
		for id := range failed {
			ids <- id
		}
		innerWg.Done()
	}()

	innerWg.Add(1)
	go func() {
		// receive from received chan
		for wm := range managers {
			dm := ms.convertToDomainManager(wm)
			err := ms.mr.Add(dm)
			if err != nil { // TODO improve error handling
				log.Println("manager service: failed to add new manager")
			}
		}
		innerWg.Done()
	}()

	workerWg.Wait()

	close(ids)
	close(failed)
	close(managers)

	innerWg.Wait()
	log.Println("manager service: AddNew returned")
	return nil
}

//
func (ms *managerService) UpdateInfos() error {
	log.Println("update infos start")
	storageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	for id := 1; id < storageManagers; id++ {
		err := ms.updateManagersInfo(id)
		if err != nil {
			log.Println("manager service:", err)
		}
	}

	log.Println("update infos stop")
	return nil
}

//
func (ms *managerService) UpdateTeams() error {
	log.Println("update teams start")
	storageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	for id := 1; id < storageManagers; id++ {
		err := ms.updateManagersTeam(id)
		if err != nil {
			log.Println("manager service:", err)
		}
	}

	log.Println("update teams stop")
	return nil
}

//
func (ms *managerService) UpdatePoints() error {
	log.Println("update points start")
	storageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	log.Println(storageManagers)

	for id := 1; id < storageManagers; id++ {
		err := ms.updateManagersPoints(id)
		if err != nil {
			log.Println("manager service:", err)
		}
	}

	log.Println("update points stop")
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

//
func (ms *managerService) updateManagersInfo(managerID int) error {
	// log.Println("update info:", managerID)
	wrapperManager, err := ms.wr.GetManager(managerID)
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
func (ms *managerService) updateManagersTeam(managerID int) error {
	// log.Println("update team:", managerID)

	gameweek, err := ms.gs.GetCurrentGameweek()
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	manager, err := ms.GetByID(managerID)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	if manager.Team.GameweekID == gameweek.ID {
		// log.Println("manager service: team already up-to-date")
		return nil
	}

	wrapperTeam, err := ms.wr.GetManagersTeam(managerID, gameweek.ID)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	team, err := ms.convertToDomainTeam(wrapperTeam)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	err = ms.mr.UpdateTeam(managerID, team)
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	return nil
}

//
func (ms *managerService) updateManagersPoints(managerID int) error {
	// log.Println("update points:", managerID)

	manager, err := ms.mr.GetByID(managerID)
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
