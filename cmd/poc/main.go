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
			fmt.Printf("!")
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
	var wg sync.WaitGroup

	for i := 0; i <= runtime.NumCPU()*16; i++ {
		go worker(&client, managerIDs, teams, &wg)
		// go fastWorker(&fastClient, managerIDs, teams)
	}

	start := time.Now()

	total := 2800
	wg.Add(total)

	go func() {
		for id := 1; id <= total; id++ {
			managerIDs <- id
		}
	}()

	tmp := make([]wrapper.Team, 0, total)

	go func() {
		for team := range teams {
			tmp = append(tmp, team)
			// fmt.Println(len(tmp))
		}
		fmt.Println("closing")
	}()

	wg.Wait()
	fmt.Println("finished after:", time.Since(start))

	close(teams)
	close(managerIDs)
	fmt.Println("closing channels")

	fmt.Println(len(tmp))
}
