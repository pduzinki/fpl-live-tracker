package manager

import (
	"fmt"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/gameweek"
	"fpl-live-tracker/pkg/storage"
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
	Update() error
	GetByID(id int) (domain.Manager, error)
}

// managerService implements ManagerService interface
type managerService struct {
	mr domain.ManagerRepository
	// ps player.PlayerService
	gs gameweek.GameweekService
	wr wrapper.Wrapper
}

// NewManagerService returns new instance of ManagerService
func NewManagerService(mr domain.ManagerRepository,
	gs gameweek.GameweekService, wr wrapper.Wrapper) (ManagerService, error) {
	rand.Seed(time.Now().UnixNano())
	ms := managerService{
		mr: mr,
		gs: gs,
		wr: wr,
	}

	return &ms, nil
}

// Update updates information about all managers currently in the game
// (i.e. updates tha data of managers already in the storage, and adds new managers)
func (ms *managerService) Update() error {
	log.Println("manager service: Update started")
	inFplManagers, err := ms.wr.GetManagersCount()
	if err != nil {
		log.Println("manager service:", err)
		return err
	}

	gameweek, err := ms.gs.GetCurrentGameweek()
	if err != nil {
		log.Println("manager service:", err)
	}

	chanSize := runtime.NumCPU() * 4
	workerCount := runtime.NumCPU() * 16

	ids := make(chan int, chanSize)
	failed := make(chan int, chanSize)
	managers := make(chan wrapper.Manager, chanSize)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	workerWg.Add(inFplManagers)

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
		for id := 1; id <= inFplManagers; id++ {
			manager, err := ms.mr.GetByID(id)
			if err != nil || manager.UpdatedInGw < gameweek.ID || gameweek.ID == 0 {
				ids <- id
			} else {
				// manager already up-to-date, skipping
				workerWg.Done()
			}
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
			err := ms.mr.Update(dm)
			if err == storage.ErrManagerNotFound {
				err := ms.mr.Add(dm)
				if err != nil {
					log.Println("manager service: failed to add new manager", err)
				}
			} else if err != nil {
				log.Println("manager service: failed to update manager", err)
			}
		}
		innerWg.Done()
	}()

	workerWg.Wait()

	close(ids)
	close(failed)
	close(managers)

	innerWg.Wait()

	log.Println("manager service: Update returned")
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
		ID:          wm.ID,
		UpdatedInGw: wm.CurrentEvent,
		Name:        fmt.Sprintf("%s %s", wm.FirstName, wm.LastName),
		TeamName:    wm.Name,
	}
}

func duration() time.Duration {
	// random duration between 30s to 5min
	return (time.Duration(rand.Intn(270)) * time.Second) + 30
}
