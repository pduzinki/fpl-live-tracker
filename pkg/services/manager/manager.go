package manager

import (
	"errors"
	"fmt"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"net/http"
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
	fplManagers, err := ms.wr.GetManagersCount()
	if err != nil {
		return err
	}

	storageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	log.Printf("fpl managers: %d	storage managers: %d\n", fplManagers, storageManagers)
	fplManagers = 42
	storageManagers = 4

	if storageManagers == fplManagers {
		return nil
	}

	var workersCount = 4
	managerIDs := make(chan int, workersCount*2)                  // to send manager ID to worker
	wrapperManagers := make(chan wrapper.Manager, workersCount*2) // to receive wrapper.Manager objects

	for worker := 1; worker <= workersCount; worker++ {
		go ms.workerAddNew(worker, managerIDs, wrapperManagers)
	}

	batchSize := 16
	for fplManagers > storageManagers {
		if fplManagers-storageManagers < batchSize {
			batchSize = fplManagers - storageManagers
		}

		for i := storageManagers + 1; i <= batchSize+storageManagers; i++ {
			managerIDs <- i
		}

		managers := make([]wrapper.Manager, 0, batchSize)
		for len(managers) != batchSize {
			managers = append(managers, <-wrapperManagers)
		}

		// TODO actually add managers to storage here
		log.Println("managers added:", batchSize)

		storageManagers += batchSize
	}
	close(managerIDs)

	log.Println("add new done")
	return nil

	// if fplManagers > storageManagers {
	// 	for id := storageManagers + 1; id <= fplManagers; id++ {
	// 		err := ms.updateManagersInfo(id)
	// 		if err != nil {
	// 			log.Println("manager service:", err)
	// 		}
	// 	}
	// }
}

//
func (ms *managerService) UpdateInfos() error {
	// storageManagers, err := ms.mr.GetCount()
	// if err != nil {
	// 	return err
	// }

	// for id := 1; id < storageManagers; id++ {
	// 	err := ms.updateManagersInfo(id)
	// 	if err != nil {
	// 		log.Println("manager service:", err)
	// 	}
	// }

	return nil
}

//
func (ms *managerService) UpdateTeams() error {
	// storageManagers, err := ms.mr.GetCount()
	// if err != nil {
	// 	return err
	// }

	// for id := 1; id < storageManagers; id++ {
	// 	err := ms.updateManagersTeam(id)
	// 	if err != nil {
	// 		log.Println("manager service:", err)
	// 	}
	// }

	return nil
}

//
func (ms *managerService) UpdatePoints() error {
	// storageManagers, err := ms.mr.GetCount()
	// if err != nil {
	// 	return err
	// }

	// log.Println(storageManagers)

	// for id := 1; id < storageManagers; id++ {
	// 	err := ms.updateManagersPoints(id)
	// 	if err != nil {
	// 		log.Println("manager service:", err)
	// 	}
	// }

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
func (ms *managerService) workerAddNew(workerID int, managerIDs chan int, wrapperManagers chan wrapper.Manager) {
	for managerID := range managerIDs {
		wrapperManager, err := ms.wr.GetManager(managerID)
		var herr *wrapper.ErrorHttpNotOk
		if errors.As(err, &herr) {
			switch herr.GetHttpStatusCode() {
			case http.StatusTooManyRequests:
				log.Printf("worker %d, 429\n", workerID)
				managerIDs <- managerID
				// TODO add sleep
				continue
			case http.StatusServiceUnavailable:
				log.Printf("worker %d, 503\n", workerID)
				managerIDs <- managerID
				// TODO add sleep
				continue
			case http.StatusNotFound:
				log.Printf("worker %d, 404\n", workerID)
				wrapperManager = wrapper.Manager{
					ID:        managerID,
					FirstName: "Not found",
					LastName:  "Not found",
					Name:      "Not found",
				}
			default:
				log.Println("other http error")
				managerIDs <- managerID
				// TODO add sleep
				continue
			}
		} else if err != nil {
			log.Println("other err", err)
			managerIDs <- managerID
			// TODO add sleep
			continue
		}

		wrapperManagers <- wrapperManager
		log.Printf("ok %d\n", wrapperManager.ID)
	}
	log.Println("worker done")
}

//
// func (ms *managerService) updateManagersInfo(managerID int) error {
// 	wrapperManager, err := ms.wr.GetManager(managerID)
// 	if err != nil {
// 		log.Println("manager service:", err)
// 		return err
// 	}

// 	manager := ms.convertToDomainManager(wrapperManager)
// 	err = ms.mr.UpdateInfo(manager.ID, manager.Info)
// 	if err == storage.ErrManagerNotFound {
// 		err = ms.mr.Add(manager)
// 		if err != nil {
// 			log.Println("manager service:", err)
// 			return err
// 		}
// 	} else if err != nil {
// 		log.Println("manager service:", err)
// 		return err
// 	}

// 	return nil
// }

//
// func (ms *managerService) updateManagersTeam(managerID int) error {
// 	gameweek, err := ms.gs.GetCurrentGameweek()
// 	if err != nil {
// 		log.Println("manager service:", err)
// 		return err
// 	}

// 	manager, err := ms.GetByID(managerID)
// 	if err != nil {
// 		log.Println("manager service:", err)
// 		return err
// 	}

// 	if manager.Team.GameweekID == gameweek.ID {
// 		// log.Println("manager service: team already up-to-date")
// 		return nil
// 	}

// 	wrapperTeam, err := ms.wr.GetManagersTeam(managerID, gameweek.ID)
// 	if err != nil {
// 		log.Println("manager service:", err)
// 		return err
// 	}

// 	team, err := ms.convertToDomainTeam(wrapperTeam)
// 	if err != nil {
// 		log.Println("manager service:", err)
// 		return err
// 	}

// 	err = ms.mr.UpdateTeam(managerID, team)
// 	if err != nil {
// 		log.Println("manager service:", err)
// 		return err
// 	}

// 	return nil
// }

//
// func (ms *managerService) updateManagersPoints(managerID int) error {
// 	manager, err := ms.mr.GetByID(managerID)
// 	if err != nil {
// 		return err
// 	}
// 	team := manager.Team

// 	err = ms.updateTeamPlayersStats(&team)
// 	if err != nil {
// 		return err
// 	}

// 	totalPoints := calculateTotalPoints(&team)
// 	subPoints := calculateSubPoints(&team)

// 	team.TotalPoints = totalPoints - team.HitPoints
// 	team.TotalPointsAfterSubs = totalPoints + subPoints - team.HitPoints

// 	// log.Println(team.TotalPoints)
// 	// log.Println(team.TotalPointsAfterSubs)

// 	err = ms.mr.UpdateTeam(manager.ID, team)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }
