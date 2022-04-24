package wrapper

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime"
	"time"
)

const DefaultURL = "https://fantasy.premierleague.com/api"

// Wrapper is a helper interface around FPL API
type Wrapper interface {
	GetClubs() ([]Club, error)
	GetFixtures() ([]Fixture, error)
	GetGameweeks() ([]Gameweek, error)
	GetPlayers() ([]Player, error)
	GetPlayersStats(gameweekID int) ([]PlayerStats, error)
	GetManagersCount() (int, error)
	GetManager(id int) (Manager, error)
	GetManagersTeam(managerID, gameweekID int) (Team, error)
}

// wrapper implements Wrapper interface
type wrapper struct {
	client  *http.Client // TODO mayber later use fasthttp
	baseURL string
}

// NewWrapper returns new instance of Wrapper.
func NewWrapper() Wrapper {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = runtime.NumCPU() * 16
	t.MaxConnsPerHost = runtime.NumCPU() * 16
	t.MaxIdleConnsPerHost = runtime.NumCPU() * 16

	return &wrapper{
		client: &http.Client{
			Timeout:   time.Second * 10,
			Transport: t,
		},
		baseURL: DefaultURL,
	}
}

// GetClubs queries https://fantasy.premierleague.com/api/bootstrap-static/
// and returns slice of wrapper.Club, or error otherwise
func (w *wrapper) GetClubs() ([]Club, error) {
	url := fmt.Sprintf(w.baseURL + "/bootstrap-static/")
	var bs Bootstrap

	err := w.fetchData(url, &bs)
	if err != nil {
		return nil, err
	}

	return bs.Clubs, nil
}

// GetFixtures queries https://fantasy.premierleague.com/api/fixtures/
// and returns slice of wrapper.Fixture, or error otherwise
func (w *wrapper) GetFixtures() ([]Fixture, error) {
	url := fmt.Sprintf(w.baseURL + "/fixtures/")
	fixtures := make([]Fixture, 380)

	err := w.fetchData(url, &fixtures)
	if err != nil {
		return nil, err
	}

	return fixtures, nil
}

// GetGameweeks queries https://fantasy.premierleague.com/api/bootstrap-static/
// and returns slice of wrapper.Gameweek, or error otherwise
func (w *wrapper) GetGameweeks() ([]Gameweek, error) {
	url := fmt.Sprintf(w.baseURL + "/bootstrap-static/")
	var bs Bootstrap

	err := w.fetchData(url, &bs)
	if err != nil {
		return nil, err
	}

	return bs.Gameweeks, nil
}

// GetPlayers queries https://fantasy.premierleague.com/api/bootstrap-static/
// and returns slice of wrapper.Player, or error otherwise
func (w *wrapper) GetPlayers() ([]Player, error) {
	url := fmt.Sprintf(w.baseURL + "/bootstrap-static/")
	var bs Bootstrap

	err := w.fetchData(url, &bs)
	if err != nil {
		return nil, err
	}

	return bs.Players, nil
}

// GetPlayersStats queries https://fantasy.premierleague.com/api/event/{gameweekID}/live/
// and returns slice of wrapper.PlayerStats, which represent live data for given gameweek
func (w *wrapper) GetPlayersStats(gameweekID int) ([]PlayerStats, error) {
	url := fmt.Sprintf(w.baseURL+"/event/%d/live/", gameweekID)
	var elements Elements

	err := w.fetchData(url, &elements)
	if err != nil {
		return nil, err
	}

	return elements.PlayersStats, nil
}

//
func (w *wrapper) GetManagersCount() (int, error) {
	// url := fmt.Sprintf(w.baseURL + "/bootstrap-static/")
	// var bs Bootstrap

	// err := w.fetchData(url, &bs)
	// if err != nil {
	// 	return 0, err
	// }

	// return bs.ManagersCount, nil

	return 100, nil // TODO remove later, let's just limit this for now
}

// GetManager queries https://fantasy.premierleague.com/api/entry/{id}/
// and returns wrapper.Manager, containing basic information on manager
func (w *wrapper) GetManager(id int) (Manager, error) {
	url := fmt.Sprintf(w.baseURL+"/entry/%d/", id)
	var m Manager

	err := w.fetchData(url, &m)
	if err != nil {
		return Manager{}, err
	}

	return m, nil
}

// GetManagersTeam queries https://fantasy.premierleague.com/api/entry/{managerID}/event/{gameweekID}/picks/
// and returns wrapper.Team, consisting of manager team picks for given gameweekf
func (w *wrapper) GetManagersTeam(managerID, gameweekID int) (Team, error) {
	url := fmt.Sprintf(w.baseURL+"/entry/%d/event/%d/picks/", managerID, gameweekID)
	var t Team

	err := w.fetchData(url, &t)
	if err != nil {
		return Team{}, err
	}

	t.ID = managerID
	return t, nil
}

// fetchData is a helper method that forms and sends http request,
// and unmarshals the response
func (w *wrapper) fetchData(url string, data interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("User-Agent", "app")

	resp, err := w.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := ErrorHttpNotOk{
			StatusCode: resp.StatusCode,
		}

		return err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ErrReadFailure
	}

	err = json.Unmarshal(body, data)
	if err != nil {
		return ErrUnmarshalFailure
	}

	return nil
}
