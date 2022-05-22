package main

import (
	"encoding/json"
	"fmt"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage/memory"
	"fpl-live-tracker/pkg/wrapper"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func worker(client *http.Client, ids <-chan int, failed chan<- int, teams chan<- wrapper.Team, wg *sync.WaitGroup) {
	for id := range ids {
		url := fmt.Sprintf("https://fantasy.premierleague.com/api/entry/%d/event/33/picks/", id)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "app")

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("err", err, "http code", resp.StatusCode)

			failed <- id
			time.Sleep(10 * time.Second)
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var team wrapper.Team

		err = json.Unmarshal(body, &team)
		if err != nil {
			log.Println("unmarshal error")
		}
		team.ID = id

		teams <- team
		wg.Done()
	}
}

func addTeams(mr domain.ManagerRepository) {
	fmt.Println("start")
	fmt.Println("core count:", runtime.NumCPU())

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = runtime.NumCPU() * 16
	t.MaxConnsPerHost = runtime.NumCPU() * 16
	t.MaxIdleConnsPerHost = runtime.NumCPU() * 16

	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	ids := make(chan int, runtime.NumCPU()*4)
	failed := make(chan int, runtime.NumCPU()*4)
	received := make(chan wrapper.Team, runtime.NumCPU()*4)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	for i := 0; i <= runtime.NumCPU()*16; i++ {
		go worker(&client, ids, failed, received, &workerWg)
	}

	start := time.Now()

	total := 2800
	workerWg.Add(total)

	innerWg.Add(1)
	go func() {
		for id := 1; id <= total; id++ {
			ids <- id
		}

		innerWg.Done()
		fmt.Println("closure 1 closing")
	}()

	innerWg.Add(1)
	go func() {
		for id := range failed {
			ids <- id
		}
		innerWg.Done()
		fmt.Println("closure 2 closing")
	}()

	innerWg.Add(1)
	go func() {
		for team := range received {
			manager := convert(team)
			err := mr.Add(manager)
			if err != nil {
				fmt.Println(err)
			}
		}
		innerWg.Done()
		fmt.Println("closure 3 closing")
	}()

	workerWg.Wait()
	fmt.Println("workers finished after:", time.Since(start))

	close(received)
	close(ids)
	close(failed)
	fmt.Println("channels closed")

	innerWg.Wait()
}

func updateTeams(mr domain.ManagerRepository) {
	fmt.Println("start")
	fmt.Println("core count:", runtime.NumCPU())

	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = runtime.NumCPU() * 16
	t.MaxConnsPerHost = runtime.NumCPU() * 16
	t.MaxIdleConnsPerHost = runtime.NumCPU() * 16

	client := http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}

	ids := make(chan int, runtime.NumCPU()*4)
	failed := make(chan int, runtime.NumCPU()*4)
	received := make(chan wrapper.Team, runtime.NumCPU()*4)
	var workerWg sync.WaitGroup
	var innerWg sync.WaitGroup

	for i := 0; i <= runtime.NumCPU()*16; i++ {
		go worker(&client, ids, failed, received, &workerWg)
	}

	start := time.Now()

	total := 2800
	workerWg.Add(total)
	alreadyUpdated := 0

	innerWg.Add(1)
	go func() {
		for id := 1; id <= total; id++ {
			// check if team with given id needs update,
			// if not (i.e. already updated),
			// then decrement workerWg,
			// else send id into channel to worker

			manager, err := mr.GetByID(id)
			if err != nil {
				fmt.Println("GetByID failed")

			}

			if manager.Team.GameweekID == 33 {
				// already updated
				alreadyUpdated++
				workerWg.Done()
				continue
			} else {
				ids <- id
			}
		}

		innerWg.Done()
		fmt.Println("closure 1 closing")
	}()

	innerWg.Add(1)
	go func() {
		for id := range failed {
			ids <- id
		}
		innerWg.Done()
		fmt.Println("closure 2 closing")
	}()

	innerWg.Add(1)
	go func() {
		for team := range received {
			manager := convert(team)
			dteam := convertTeam(team)
			err := mr.UpdateTeam(manager.ID, dteam)
			if err != nil {
				fmt.Println(err)
			}
		}
		innerWg.Done()
		fmt.Println("closure 3 closing")
	}()

	workerWg.Wait()
	fmt.Println("workers finished after:", time.Since(start))

	close(received)
	close(ids)
	close(failed)
	fmt.Println("channels closed")

	innerWg.Wait()
	fmt.Println("team already updated:", alreadyUpdated)
}

func main() {
	mr := memory.NewManagerRepository()
	addTeams(mr)
	added, _ := mr.GetCount()
	fmt.Println("new managers added:", added)

	// now that managers are added, let's try to update them, but only if they really need update (i.e GameweekID < live gameweek)
	fmt.Println("----")
	fmt.Println("updating teams...")
	updateTeams(mr)
}

func convert(wt wrapper.Team) domain.Manager {
	manager := domain.Manager{
		ID: wt.ID,
		Team: domain.Team{
			GameweekID: wt.EntryHistory.GameweekID,
			Picks:      nil,
			ActiveChip: wt.ActiveChip,
			HitPoints:  wt.EntryHistory.EventTransfersCost,
		},
	}

	return manager
}

func convertTeam(wt wrapper.Team) domain.Team {
	return domain.Team{
		GameweekID: wt.EntryHistory.GameweekID,
		Picks:      nil,
		ActiveChip: wt.ActiveChip,
		HitPoints:  wt.EntryHistory.EventTransfersCost,
	}
}
