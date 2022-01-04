package wrapper

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
	"time"

	domain "fpl-live-tracker/pkg"
)

func TestGetManager(t *testing.T) {
	testcases := []struct {
		name                string
		id                  int
		want                *domain.Manager
		wantErr             error
		handlerStatusCode   int
		handlerBodyFilePath string
	}{
		{"ok", 123, &domain.Manager{FplID: 123, FullName: "John Doe", TeamName: "FC John"}, nil, http.StatusOK, "./testdata/getmanager.json"},
		{"too many requests", 123, nil, errorHttpNotOk{429}, http.StatusTooManyRequests, "./testdata/empty.json"},
		{"not found", 123, nil, errorHttpNotOk{404}, http.StatusNotFound, "./testdata/empty.json"},
		{"service unavailable", 123, nil, errorHttpNotOk{503}, http.StatusServiceUnavailable, "./testdata/empty.json"},
		{"unmarshal error", 123, nil, ErrUnmarshalFailure, http.StatusOK, "./testdata/getmanager.broken_json"},
	}

	for _, test := range testcases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.handlerStatusCode)
			w.Header().Set("Content-Type", "application/json")

			f, err := os.ReadFile(test.handlerBodyFilePath)
			if err != nil {
				t.Error(err)
			}

			w.Write(f)
		}))
		defer server.Close()

		w := NewWrapper(server.URL)

		got, gotErr := w.GetManager(test.id)
		// if gotErr != test.wantErr {
		if !errors.Is(gotErr, test.wantErr) {
			t.Errorf("error: testcase '%s', want error '%v', got error '%v'", test.name, test.wantErr, gotErr)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: testcase '%s', want '%v', got '%v'", test.name, test.want, got)
		}
	}
}

func TestGetGameweeks(t *testing.T) {
	dt, _ := time.Parse(time.RFC3339, "2021-11-20T11:00:00Z")

	testcases := []struct {
		name                string
		handlerStatusCode   int
		handlerBodyFilePath string
		wantErr             error
		want                domain.Gameweek // verify just one particular gw data, just to save some space
	}{
		{"ok", http.StatusOK, "./testdata/bootstrap-static.json", nil, domain.Gameweek{
			ID:           12,
			Name:         "Gameweek 12",
			Finished:     false,
			IsCurrent:    true,
			IsNext:       false,
			DeadlineTime: dt,
		}},
	}

	for _, test := range testcases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.handlerStatusCode)
			w.Header().Set("Content-Type", "application/json")

			f, err := os.ReadFile(test.handlerBodyFilePath)
			if err != nil {
				t.Error(err)
			}

			w.Write(f)
		}))
		defer server.Close()

		w := NewWrapper(server.URL)

		got, err := w.GetGameweeks()
		if err != test.wantErr {
			t.Errorf("error: testcase '%s', want error '%v', got error '%v'", test.name, test.wantErr, err)
		}
		if got[11] != test.want {
			t.Errorf("error: testcase '%s', want '%v', got '%v'", test.name, test.want, got[11])
		}
	}
}

/*
func TestGetFixtures(t *testing.T) {
	testcases := []struct {
		name                string
		handlerStatusCode   int
		handlerBodyFilePath string
		wantErr             error
		want                domain.Fixture
	}{
		{
			name:                "ok",
			handlerStatusCode:   200,
			handlerBodyFilePath: "./testdata/fixtures17.json",
			wantErr:             nil,
			want: domain.Fixture{
				ID: 163,
			},
		},
	}

	for _, test := range testcases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.handlerStatusCode)
			w.Header().Set("Content-Type", "application/json")

			f, err := os.ReadFile(test.handlerBodyFilePath)
			if err != nil {
				t.Error(err)
			}

			w.Write(f)
		}))
		defer server.Close()

		w := NewWrapper(server.URL)

		got, err := w.GetFixtures(17)
		if err != test.wantErr {
			t.Errorf("error: testcase '%s', want error '%v', got error '%v'", test.name, test.wantErr, err)
		}
		if got[2] != test.want {
			t.Errorf("error: testcase '%s', want '%v', got '%v'", test.name, test.want, got[2])
		}
	}
}
TODO fix this later
*/

func TestGetFixtures(t *testing.T) {
	w := NewWrapper(DefaultURL)

	f, _ := w.GetFixtures()
	_ = f

	i := 0
	i++
}

func TestFetchData(t *testing.T) {
	testcases := []struct {
		name                string
		handlerStatusCode   int
		handlerBodyFilePath string
		wantErr             error
	}{
		{"ok", http.StatusOK, "./testdata/fetchdata.json", nil},
		{"too many requests", http.StatusTooManyRequests, "./testdata/fetchdata.json", errorHttpNotOk{429}},
		{"not found", http.StatusNotFound, "./testdata/fetchdata.json", errorHttpNotOk{404}},
		{"service unavailable", http.StatusServiceUnavailable, "./testdata/fetchdata.json", errorHttpNotOk{503}},
		{"unmarshal error", http.StatusOK, "./testdata/fetchdata.broken_json", ErrUnmarshalFailure},
	}

	for _, test := range testcases {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(test.handlerStatusCode)
			w.Header().Set("Content-Type", "application/json")

			f, err := os.ReadFile(test.handlerBodyFilePath)
			if err != nil {
				t.Error(err)
			}

			w.Write(f)
		}))
		defer server.Close()

		w := wrapper{
			client:  &http.Client{},
			baseURL: server.URL,
		}

		type tmp struct {
			Data int `json:"data"`
		}
		var data tmp

		gotErr := w.fetchData(w.baseURL, &data)
		if gotErr != test.wantErr {
			t.Errorf("error: testcase '%s', want error '%v', got error '%v'", test.name, test.wantErr, gotErr)
		}
	}
}
