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

	teams := make([]wrapper.Team, 0, total)

	innerWg.Add(1)
	go func() {
		for team := range received {
			teams = append(teams, team)
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
	fmt.Println(len(teams))
}
