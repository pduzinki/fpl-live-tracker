package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
)

var (
	john = domain.Manager{ID: 1, Name: "John Doe", TeamName: "FC John"}
	jim  = domain.Manager{ID: 2, Name: "Jim Jim", TeamName: "FC Jim"}
	jane = domain.Manager{ID: 3, Name: "Jane Foo", TeamName: "Jane City"}
	joel = domain.Manager{ID: 66, Name: "Joel Bar", TeamName: "Bar AFC"}
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
			john.ID: john,
			jim.ID:  jim,
		},
	}

	for _, test := range testcases {
		got := mr.Add(test.manager)
		if got != test.want {
			t.Errorf("error: for %v, got err '%v', want '%v'", test.manager, got, test.want)
		}

		if v, ok := mr.managers[test.manager.ID]; ok {
			if !reflect.DeepEqual(v, test.manager) {
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
			john.ID: john,
			jim.ID:  jim,
		},
	}

	for _, test := range testcases {
		got := mr.AddMany(test.managers)
		if got != test.want {
			t.Errorf("error: got err '%v', want err '%v'", got, test.want)
		}

		for _, m := range test.managers {
			if v, ok := mr.managers[m.ID]; ok {
				if !reflect.DeepEqual(v, m) {
					t.Errorf("error: incorrect manager data in memory storage")
				}
			} else {
				t.Errorf("error: manager not found in memory storage")
			}

		}
	}
}

func TestManagerUpdate(t *testing.T) {
	testcases := []struct {
		manager domain.Manager
		want    error
	}{
		{
			manager: domain.Manager{
				ID:       john.ID,
				Name:     john.Name,
				TeamName: "John United",
			},
			want: nil,
		},
		{
			manager: domain.Manager{
				ID: 0,
			},
			want: storage.ErrManagerNotFound,
		},
	}

	mr := managerRepository{
		managers: map[int]domain.Manager{
			john.ID: john,
		},
	}

	for _, test := range testcases {
		got := mr.Update(test.manager)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}

		if got == nil {
			if v, ok := mr.managers[test.manager.ID]; ok {
				if !reflect.DeepEqual(v, test.manager) {
					t.Errorf("error: incorrect manager data in memory storage")
				}
			} else {
				t.Errorf("error: manager not found in memory storage")
			}
		}
	}
}

func TestManagerGetByID(t *testing.T) {
	testcases := []struct {
		fplID   int
		want    domain.Manager
		wantErr error
	}{
		{john.ID, john, nil},
		{jane.ID, domain.Manager{}, storage.ErrManagerNotFound},
	}

	mr := managerRepository{
		managers: map[int]domain.Manager{
			john.ID: john,
			jim.ID:  jim,
		},
	}

	for _, test := range testcases {
		got, gotErr := mr.GetByID(test.fplID)
		if gotErr != test.wantErr {
			t.Errorf("error: for %v, got err '%v', want err '%v'", test.fplID, gotErr, test.wantErr)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: for %v, got '%v', want '%v'", test.fplID, got, test.want)
		}
	}
}

func TestManagerGetCount(t *testing.T) {
	testcases := []struct {
		want    int
		wantErr error
	}{
		{
			want:    2,
			wantErr: nil,
		},
	}

	mr := managerRepository{
		managers: map[int]domain.Manager{
			john.ID: john,
			jim.ID:  jim,
		},
	}

	for _, test := range testcases {
		got, gotErr := mr.GetCount()
		if gotErr != test.wantErr {
			t.Errorf("error: got err '%v', want err '%v'", gotErr, test.wantErr)
		}
		if got != test.want {
			t.Errorf("error: got '%v', want '%v'", got, test.want)
		}
	}
}
