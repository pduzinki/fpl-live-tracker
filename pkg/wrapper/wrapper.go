package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	domain "fpl-live-tracker/pkg"
)

const DefaultURL = "https://fantasy.premierleague.com/api"

var ErrReadFailure error = errors.New("failed to read the response")
var ErrUnmarshalFailure error = errors.New("failed to unmarshal data")

type errorHttpNotOk struct {
	statusCode int
}

type ErrorHttpNotOk interface {
	error
	GetHttpStatusCode() int
}

func (err errorHttpNotOk) Error() string {
	return fmt.Sprintf("http status not ok: %d\n", err.statusCode)
}

func (err errorHttpNotOk) GetHttpStatusCode() int {
	return err.statusCode
}

// Wrapper is a helper interface around FPL API
type Wrapper interface {
	GetManager(id int) (*domain.Manager, error)
	GetTeam(id, gw int) (*domain.Team, error)
	GetGameweeks() ([]domain.Gameweek, error)
}

type wrapper struct {
	client  *http.Client // TODO mayber later use fasthttp
	baseURL string
}

// NewWrapper returns an instance of an FPL API wrapper.
// Pass wrapper.DefaultURL as an argument, if you're not testing anything.
func NewWrapper(url string) Wrapper {
	return &wrapper{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: url,
	}
}

// GetManager returns data from FPL API "/api/entry/{managerID}/" endpoint
func (w *wrapper) GetManager(id int) (*domain.Manager, error) {
	url := fmt.Sprintf(w.baseURL+"/entry/%d/", id)
	var m Manager

	err := w.fetchData(url, &m)
	if err != nil {
		return nil, err
	}

	tm := domain.Manager{
		FplID:    m.ID,
		FullName: fmt.Sprintf("%s %s", m.FirstName, m.LastName),
		TeamName: m.Name,
	}

	return &tm, nil
}

//
func (w *wrapper) GetTeam(id, gw int) (*domain.Team, error) {
	url := fmt.Sprintf(w.baseURL+"/entry/%d/event/%d/picks/", id, gw)
	var t Team

	err := w.fetchData(url, &t)
	if err != nil {
		return nil, err
	}

	// TODO add proper translation wrapper.Team -> tracker.Team
	tt := domain.Team{
		FplID: id,
	}

	return &tt, nil
}

func (w *wrapper) GetGameweeks() ([]domain.Gameweek, error) {
	url := fmt.Sprintf(w.baseURL + "/bootstrap-static/")
	var bs Bootstrap

	err := w.fetchData(url, &bs)
	if err != nil {
		return nil, err
	}

	gameweeks := make([]domain.Gameweek, 0)
	for _, gw := range bs.Gws {
		deadlineTime, err := time.Parse(time.RFC3339, gw.DeadlineTime)
		if err != nil {
			panic(err)
		}

		gameweeks = append(gameweeks, domain.Gameweek{
			ID:           gw.ID,
			Name:         gw.Name,
			Finished:     gw.Finished,
			IsCurrent:    gw.IsCurrent,
			IsNext:       gw.IsNext,
			DeadlineTime: deadlineTime,
		})
	}

	return gameweeks, nil
}

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
		err := errorHttpNotOk{
			statusCode: resp.StatusCode,
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
