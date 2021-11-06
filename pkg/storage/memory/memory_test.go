package memory

import (
	tracker "fpl-live-tracker/pkg"
	"testing"
)

var (
	john = tracker.Manager{FplID: 1, FullName: "John Doe", TeamName: "FC John"}
	jim  = tracker.Manager{FplID: 2, FullName: "Jim Jim", TeamName: "FC Jim"}
	jane = tracker.Manager{FplID: 3, FullName: "Jane Foo", TeamName: "Jane City"}
	joel = tracker.Manager{FplID: 66, FullName: "Joel Bar", TeamName: "Bar AFC"}
)

func TestAdd(t *testing.T) {
	testcases := []struct {
		manager tracker.Manager
		want    error
	}{
		{jane, nil},
		{john, ErrRecordAlreadyExists},
	}

	mr := managerRepository{
		managers: map[int]tracker.Manager{
			john.FplID: john,
			jim.FplID:  jim,
		},
	}

	for _, test := range testcases {
		got := mr.Add(test.manager)
		if got != test.want {
			t.Errorf("error: for %v, got '%v', want '%v'", test.manager, got, test.want)
		}
	}
}

func TestAddMany(t *testing.T) {
	// TODO
}

func TestGetByFplID(t *testing.T) {
	testcases := []struct {
		fplID   int
		want    tracker.Manager
		wantErr error
	}{
		{john.FplID, john, nil},
		{jane.FplID, tracker.Manager{}, ErrRecordNotFound},
	}

	mr := managerRepository{
		managers: map[int]tracker.Manager{
			john.FplID: john,
			jim.FplID:  jim,
		},
	}

	for _, test := range testcases {
		got, gotErr := mr.GetByFplID(test.fplID)
		if gotErr != test.wantErr {
			t.Errorf("error: for %v, got err '%v', want err '%v'", test.fplID, gotErr, test.wantErr)
		}
		if got != test.want {
			t.Errorf("error: for %v, got '%v', want '%v'", test.fplID, got, test.want)
		}
	}
}
