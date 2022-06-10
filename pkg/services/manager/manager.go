package manager

import (
	"fmt"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/services/player"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"math/rand"
	"net/http"
	"runtime"
	"sync"
	"time"
)

// ManagerService is an interface for interacting with managers
type ManagerService interface {
	AddNew() error
	UpdateInfos() error
	UpdateTeams() error
	UpdatePoints() error
	GetByID(id int) (domain.Manager, error)
}

// managerService implements ManagerService interface
type managerService struct {
	mr domain.ManagerRepository
	ps player.PlayerService
	gs gameweek.GameweekService
	wr wrapper.Wrapper
}

// NewManagerService returns new instance of ManagerService
func NewManagerService(mr domain.ManagerRepository, ps player.PlayerService,
	gs gameweek.GameweekService, wr wrapper.Wrapper) (ManagerService, error) {
	rand.Seed(time.Now().UnixNano())
	ms := managerService{
		mr: mr,
		ps: ps,
		gs: gs,
		wr: wr,
	}

	return &ms, nil
}

// AddNew adds managers that joined the game since the last AddNew call, to the storage
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

	for i := 0; i < workerCount; i++ {
		go func() {
			for id := range ids {
				wm, err := ms.wr.GetManager(id)
				if herr, ok := err.(wrapper.ErrorHttpNotOk); ok {
					statusCode := herr.GetHttpStatusCode()
					switch statusCode {
					case http.StatusTooManyRequests:
						log.Println("manager service: too many requests!")
						failed <- id
						time.Sleep(duration())
						continue
					case http.StatusServiceUnavailable:
						failed <- id
						time.Sleep(10 * time.Minute)
						continue
					case http.StatusNotFound:
						wm = wrapper.Manager{
							ID:        id,
							FirstName: "not found",
							LastName:  "not found",
							Name:      "not found",
						}
					default:
						failed <- id
						time.Sleep(10 * time.Minute)
						continue
					}
				} else if err != nil {
					failed <- id
					time.Sleep(10 * time.Minute)
					continue
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
		// receive from managers chan
		for wm := range managers {
			dm := ms.convertToDomainManager(wm)
			err := ms.mr.Add(dm)
			if err != nil {
				log.Println("manager service: failed to add new manager", err)
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

// UpdateInfos updates information about all managers currently in the storage
func (ms *managerService) UpdateInfos() error {
	log.Println("manager service: UpdateInfos started")
	inStorageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	chanSize := runtime.NumCPU() * 4
	workerCount := runtime.NumCPU() * 16

	ids := make(chan int, chanSize)
	failed := make(chan int, chanSize)
	managers := make(chan wrapper.Manager, chanSize)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	workerWg.Add(inStorageManagers)

	for i := 0; i < workerCount; i++ {
		go func() {
			for id := range ids {
				wm, err := ms.wr.GetManager(id)
				if herr, ok := err.(wrapper.ErrorHttpNotOk); ok {
					statusCode := herr.GetHttpStatusCode()
					switch statusCode {
					case http.StatusTooManyRequests:
						log.Println("manager service: too many requests!")
						failed <- id
						time.Sleep(duration())
						continue
					case http.StatusServiceUnavailable:
						failed <- id
						time.Sleep(10 * time.Minute)
						continue
					case http.StatusNotFound:
						wm = wrapper.Manager{
							ID:        id,
							FirstName: "not found",
							LastName:  "not found",
							Name:      "not found",
						}
					default:
						failed <- id
						time.Sleep(10 * time.Minute)
						continue
					}
				} else if err != nil {
					failed <- id
					time.Sleep(10 * time.Minute)
					continue
				}

				managers <- wm
				workerWg.Done()
			}
		}()
	}

	innerWg.Add(1)
	go func() {
		// send to ids chan
		for id := 1; id <= inStorageManagers; id++ {
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
		// receive from managers chan
		for wm := range managers {
			dm := ms.convertToDomainManager(wm)
			err := ms.mr.UpdateInfo(dm.ID, dm.Info)
			if err != nil {
				log.Println("manager service: failed to add new manager", err)
			}
		}
		innerWg.Done()
	}()

	workerWg.Wait()

	close(ids)
	close(failed)
	close(managers)

	innerWg.Wait()

	log.Println("manager service: UpdateInfos returned")
	return nil
}

// UpdateTeams updates team information for all managers currently in the storage
func (ms *managerService) UpdateTeams() error {
	// log.Println("manager service: UpdateTeams started")
	// inStorageManagers, err := ms.mr.GetCount()
	// if err != nil {
	// 	return err
	// }

	// gameweek, err := ms.gs.GetCurrentGameweek()
	// if err != nil {
	// 	log.Println("manager service:", err)
	// 	return err
	// }

	// chanSize := runtime.NumCPU() * 4
	// workerCount := runtime.NumCPU() * 16

	// ids := make(chan int, chanSize)
	// failed := make(chan int, chanSize)
	// teams := make(chan wrapper.Team, chanSize)
	// var workerWg sync.WaitGroup
	// var innerWg sync.WaitGroup

	// workerWg.Add(inStorageManagers)

	// for i := 0; i < workerCount; i++ {
	// 	go func() {
	// 		for id := range ids {
	// 			wt, err := ms.wr.GetManagersTeam(id, gameweek.ID)

	// 			if herr, ok := err.(wrapper.ErrorHttpNotOk); ok {
	// 				statusCode := herr.GetHttpStatusCode()
	// 				switch statusCode {
	// 				case http.StatusTooManyRequests:
	// 					log.Println("manager service: too many requests!")
	// 					failed <- id
	// 					time.Sleep(duration())
	// 					continue
	// 				case http.StatusServiceUnavailable:
	// 					failed <- id
	// 					time.Sleep(10 * time.Minute)
	// 					continue
	// 				case http.StatusNotFound:
	// 					wt = wrapper.Team{
	// 						ID:         id,
	// 						ActiveChip: "not found",
	// 					}
	// 				default:
	// 					failed <- id
	// 					time.Sleep(10 * time.Minute)
	// 					continue
	// 				}
	// 			} else if err != nil {
	// 				failed <- id
	// 				time.Sleep(10 * time.Minute)
	// 				continue
	// 			}

	// 			teams <- wt
	// 			workerWg.Done()
	// 		}
	// 	}()
	// }

	// innerWg.Add(1)
	// go func() {
	// 	// send to ids chan
	// 	for id := 1; id <= inStorageManagers; id++ {
	// 		ids <- id
	// 	}
	// 	innerWg.Done()
	// }()

	// innerWg.Add(1)
	// go func() {
	// 	// receive from failed chan, send to ids chan
	// 	for id := range failed {
	// 		ids <- id
	// 	}
	// 	innerWg.Done()
	// }()

	// innerWg.Add(1)
	// go func() {
	// 	// receive from managers chan
	// 	for wt := range teams {
	// 		dt, err := ms.convertToDomainTeam(wt)
	// 		if err != nil {
	// 			log.Println("manager service: failed to convert team data")
	// 		}

	// 		err = ms.mr.UpdateTeam(wt.ID, dt)
	// 		if err != nil {
	// 			log.Println("manager service: failed to add new manager", err)
	// 		}
	// 	}
	// 	innerWg.Done()
	// }()

	// workerWg.Wait()

	// close(ids)
	// close(failed)
	// close(teams)

	// innerWg.Wait()

	// log.Println("manager service: UpdateTeams returned")
	return nil
}

// UpdatePoints updates points data for all managers currently in the storage
func (ms *managerService) UpdatePoints() error {
	log.Println("manager service: UpdatePoints started")

	inStorageManagers, err := ms.mr.GetCount()
	if err != nil {
		return err
	}

	chanSize := runtime.NumCPU()
	workerCount := runtime.NumCPU()

	ids := make(chan int, chanSize)
	failed := make(chan int, chanSize)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	workerWg.Add(inStorageManagers)
	for i := 0; i < workerCount; i++ {
		go func() {
			for id := range ids {
				err := ms.updateManagersPoints(id)
				if err != nil {
					failed <- id
					log.Println("manager service: failed to update points", err)
					continue
				}
				workerWg.Done()
			}
		}()
	}

	innerWg.Add(1)
	go func() {
		// send to ids chan
		for id := 1; id <= inStorageManagers; id++ {
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

	workerWg.Wait()

	close(ids)
	close(failed)

	innerWg.Wait()

	log.Println("managers service: UpdatePoints returned")
	return nil
}

// GetByID returns managers with given ID, or error otherwise
func (ms *managerService) GetByID(id int) (domain.Manager, error) {
	manager := domain.Manager{ID: id}

	err := runManagerValidations(&manager, idHigherThanZero)
	if err != nil {
		return domain.Manager{}, err
	}

	return ms.mr.GetByID(id)
}

// convertToDomainManager returns domain.Manager, consistent with given wrapper.Manager
func (ms *managerService) convertToDomainManager(wm wrapper.Manager) domain.Manager {
	return domain.Manager{
		ID: wm.ID,
		Info: domain.ManagerInfo{
			Name:     fmt.Sprintf("%s %s", wm.FirstName, wm.LastName),
			TeamName: wm.Name,
		},
	}
}

// convertToDomainTeam returns domain.Team, consistent with given wrapper.Team
// func (ms *managerService) convertToDomainTeam(wt wrapper.Team) (domain.Team, error) {
// 	team := domain.Team{
// 		GameweekID: wt.EntryHistory.GameweekID,
// 		Picks:      make([]domain.TeamPlayer, 0, 15),
// 		ActiveChip: wt.ActiveChip,
// 		HitPoints:  wt.EntryHistory.EventTransfersCost,
// 	}

// 	for _, pick := range wt.Picks {
// 		p, err := ms.ps.GetByID(pick.ID)
// 		if err != nil {
// 			log.Println("manager service:", err)
// 			return domain.Team{}, err
// 		}

// 		dp := domain.TeamPlayer{
// 			Player:        p,
// 			IsCaptain:     pick.IsCaptain,
// 			IsViceCaptain: pick.IsViceCaptain,
// 		}

// 		team.Picks = append(team.Picks, dp)
// 	}

// 	return team, nil
// }

// updateTeamPlayersStats updates players stats in the given team
// func (ms *managerService) updateTeamPlayersStats(team *domain.Team) error {
// 	for i := 0; i < len(team.Picks); i++ {
// 		tp := team.Picks[i]
// 		p, err := ms.ps.GetByID(tp.ID)
// 		if err != nil {
// 			log.Println("manager service: failed to update manager's team stats", err)
// 			return err
// 		}
// 		tp.Stats = p.Stats
// 		team.Picks[i] = tp
// 	}

// 	return nil
// }

// updateManagersPoints updates points gained by manager's team with given ID
func (ms *managerService) updateManagersPoints(managerID int) error {
	// manager, err := ms.mr.GetByID(managerID)
	// if err != nil {
	// 	return err
	// }
	// team := manager.Team

	// err = ms.updateTeamPlayersStats(&team)
	// if err != nil {
	// 	return err
	// }

	// totalPoints := calculateTotalPoints(&team)
	// subPoints := calculateSubPoints(&team)

	// team.TotalPoints = totalPoints - team.HitPoints
	// team.TotalPointsAfterSubs = totalPoints + subPoints - team.HitPoints

	// err = ms.mr.UpdateTeam(manager.ID, team)
	// if err != nil {
	// 	return err
	// }

	return nil
}

func duration() time.Duration {
	// random duration between 30s to 5min
	return (time.Duration(rand.Intn(270)) * time.Second) + 30
}
