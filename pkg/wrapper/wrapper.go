package wrapper

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"time"

	tracker "fpl-live-tracker/pkg"
)

const baseURL = "https://fantasy.premierleague.com/api"

var ErrHTTPStatusNotOK error = errors.New("Response status != http 200")
var ErrHTTPTooManyRequests error = errors.New("Response status - http 429 - too many requests")
var ErrHTTPStatusNotFound error = errors.New("Response status - http 404 - not found")
var ErrHTTPStatusServiceUnavailable error = errors.New("Response status - http 503 - service unavailable")
var ErrReadFailure error = errors.New("Failed to read the response")
var ErrUnmarshalFailure error = errors.New("Failed to unmarshal data")

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
func NewWrapper() Wrapper {
	return &wrapper{
		client: &http.Client{
			Timeout: time.Second * 10,
		},
		baseURL: baseURL,
	}
}

func (w *wrapper) GetManager(id int) (*tracker.Manager, error) {
	return nil, nil
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

	if resp.StatusCode == http.StatusTooManyRequests {
		return ErrHTTPTooManyRequests
	} else if resp.StatusCode == http.StatusServiceUnavailable {
		return ErrHTTPStatusServiceUnavailable
	} else if resp.StatusCode == http.StatusNotFound {
		return ErrHTTPStatusNotFound
	} else if resp.StatusCode != http.StatusOK {
		return ErrHTTPStatusNotOK
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
