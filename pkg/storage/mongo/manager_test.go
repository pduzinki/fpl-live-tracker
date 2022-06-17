package mongo

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
)

var mr domain.ManagerRepository

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

	for _, test := range testcases {
		got := mr.Add(test.manager)
		if got != test.want {
			t.Errorf("error: for %v, got err '%v', want '%v'", test.manager, got, test.want)
		}
	}
}

func TestManagerAddMany(t *testing.T) {
	// TODO add test
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

	for _, test := range testcases {
		got := mr.Update(test.manager)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}
	}
}

func TestManagerGetByID(t *testing.T) {
	testcases := []struct {
		fplID   int
		want    domain.Manager
		wantErr error
	}{
		{jane.ID, jane, nil},
		{joel.ID, domain.Manager{}, storage.ErrManagerNotFound},
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
