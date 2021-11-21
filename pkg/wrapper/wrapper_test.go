package wrapper

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	tracker "fpl-live-tracker/pkg"
)

func TestGetManager(t *testing.T) {
	testcases := []struct {
		name                string
		id                  int
		want                *tracker.Manager
		wantErr             error
		handlerStatusCode   int
		handlerBodyFilePath string
	}{
		{"ok", 123, &tracker.Manager{FplID: 123, FullName: "John Doe", TeamName: "FC John"}, nil, http.StatusOK, "./testdata/getmanager.json"},
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
	testcases := []struct {
		name                string
		handlerStatusCode   int
		handlerBodyFilePath string
	}{
		{"ok", http.StatusOK, "./testdata/bootstrap-static.json"},
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

		_, err := w.GetGameweeks()
		fmt.Println(err)

		// TODO finish up this test
	}
}

func TestFetchData(t *testing.T) {
	// TODO add this test

	testcases := []struct {
		foo int
	}{
		{0},
		{1},
		{2},
	}

	for _, test := range testcases {
		_ = test
	}
}
