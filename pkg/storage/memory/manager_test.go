package memory

import (
	domain "fpl-live-tracker/pkg"
	"testing"
)

var (
	john = domain.Manager{FplID: 1, FullName: "John Doe", TeamName: "FC John"}
	jim  = domain.Manager{FplID: 2, FullName: "Jim Jim", TeamName: "FC Jim"}
	jane = domain.Manager{FplID: 3, FullName: "Jane Foo", TeamName: "Jane City"}
	joel = domain.Manager{FplID: 66, FullName: "Joel Bar", TeamName: "Bar AFC"}
)

func TestManagerAdd(t *testing.T) {
	testcases := []struct {
		manager domain.Manager
		want    error
	}{
		{jane, nil},
		{john, errManagerAlreadyExists{john.FplID}},
	}

	mr := managerRepository{
		managers: map[int]domain.Manager{
			john.FplID: john,
			jim.FplID:  jim,
		},
	}

	for _, test := range testcases {
		got := mr.Add(test.manager)
		if got != test.want {
			t.Errorf("error: for %v, got '%v', want '%v'", test.manager, got, test.want)
		}
		// TODO check if manager was really saved in storage
	}
}

func TestManagerAddMany(t *testing.T) {
	testcases := []struct {
		managers []domain.Manager
		want     error
	}{
		{
			managers: []domain.Manager{jane, joel},
			want:     nil,
		},
		{
			managers: []domain.Manager{jim, jim},
			want:     errManagerAlreadyExists{jim.FplID},
		},
	}

	mr := managerRepository{
		managers: map[int]domain.Manager{
			john.FplID: john,
			jim.FplID:  jim,
		},
	}

	for _, test := range testcases {
		got := mr.AddMany(test.managers)
		if got != test.want {
			t.Errorf("error: got err '%v', want err '%v'", got, test.want)
		}

		// TODO check if managers were really saved in storage
	}
}

func TestManagerGetByFplID(t *testing.T) {
	testcases := []struct {
		fplID   int
		want    domain.Manager
		wantErr error
	}{
		{john.FplID, john, nil},
		{jane.FplID, domain.Manager{}, ErrManagerNotFound},
	}

	mr := managerRepository{
		managers: map[int]domain.Manager{
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
