package main

import (
	"encoding/json"
	"fmt"
	"fpl-live-tracker/pkg/wrapper"
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"sync"
	"time"
)

func worker(client *http.Client, managerIDs chan int, teams chan wrapper.Team, wg *sync.WaitGroup) {
	for id := range managerIDs {
		url := fmt.Sprintf("https://fantasy.premierleague.com/api/entry/%d/event/33/picks/", id)

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "app")

		resp, err := client.Do(req)
		if err != nil || resp.StatusCode != http.StatusOK {
			fmt.Println("err", err, "http code", resp.StatusCode)
			teams <- wrapper.Team{}
			wg.Done()
			continue
		}

		body, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		var team wrapper.Team

		err = json.Unmarshal(body, &team)
		if err != nil {
			log.Println("unmarshal error")
		}

		teams <- team
		wg.Done()
	}
}

func main() {
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

	managerIDs := make(chan int, runtime.NumCPU()*4)
	teams := make(chan wrapper.Team, runtime.NumCPU()*4)
	var workerWg sync.WaitGroup
	var closureWg sync.WaitGroup

	for i := 0; i <= runtime.NumCPU()*8; i++ {
		go worker(&client, managerIDs, teams, &workerWg)
	}

	start := time.Now()

	total := 2800
	workerWg.Add(total)
	closureWg.Add(2)

	go func() {
		for id := 1; id <= total; id++ {
			managerIDs <- id
		}
		closureWg.Done()
		fmt.Println("closure 1 closing")
	}()

	tmp := make([]wrapper.Team, 0, total)

	go func() {
		for team := range teams {
			tmp = append(tmp, team)
		}
		closureWg.Done()
		fmt.Println("closure 2 closing")
	}()

	workerWg.Wait()
	fmt.Println("workers finished after:", time.Since(start))

	close(teams)
	close(managerIDs)
	fmt.Println("channels closed")

	closureWg.Wait()
	fmt.Println(len(tmp))
}
