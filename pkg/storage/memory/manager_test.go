package memory

import (
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/storage"
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
		{john, storage.ErrManagerAlreadyExists},
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
			t.Errorf("error: for %v, got err '%v', want '%v'", test.manager, got, test.want)
		}

		if v, ok := mr.managers[test.manager.FplID]; ok {
			if v != test.manager {
				t.Errorf("error: incorrect manager data in memory storage")
			}
		} else {
			t.Errorf("error: manager not found in memory storage")
		}
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
			want:     storage.ErrManagerAlreadyExists,
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

		for _, m := range test.managers {
			if v, ok := mr.managers[m.FplID]; ok {
				if v != m {
					t.Errorf("error: incorrect manager data in memory storage")
				}
			} else {
				t.Errorf("error: manager not found in memory storage")
			}

		}
	}
}

func TestManagerGetByFplID(t *testing.T) {
	testcases := []struct {
		fplID   int
		want    domain.Manager
		wantErr error
	}{
		{john.FplID, john, nil},
		{jane.FplID, domain.Manager{}, storage.ErrManagerNotFound},
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
