package wrapper

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

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
