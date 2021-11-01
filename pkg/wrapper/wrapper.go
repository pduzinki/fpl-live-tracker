package wrapper

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	tracker "fpl-live-tracker/pkg"
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

func (err *errorHttpNotOk) Error() string {
	return fmt.Sprintf("http status not ok: %d\n", err.statusCode)
}

func (err *errorHttpNotOk) GetHttpStatusCode() int {
	return err.statusCode
}

// Wrapper is a helper interface around FPL API
type Wrapper interface {
	GetManager(id int) (*tracker.Manager, error)
	// GetTeam(id int) (*tracker.Team, error)
	// TODO add more methods
}

type wrapper struct {
	client  *http.Client // TODO mayber later use fasthttp
	baseURL string
}

// NewWrapper returns an instance of an FPL API wrapper
func NewWrapper(url string) Wrapper {
	return &wrapper{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: url,
	}
}

func (w *wrapper) GetManager(id int) (*tracker.Manager, error) {
	url := fmt.Sprintf(w.baseURL+"/entry/%d/", id)
	var m Manager

	err := w.fetchData(url, &m)
	if err != nil {
		return nil, err
	}

	tm := tracker.Manager{
		FplID:    m.ID,
		FullName: fmt.Sprintf("%s %s", m.FirstName, m.LastName),
		TeamName: m.Name,
	}

	return &tm, nil
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

		return &err
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
