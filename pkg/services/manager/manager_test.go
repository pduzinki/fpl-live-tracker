package manager

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
)

var john = domain.Manager{ID: 1, Info: domain.ManagerInfo{Name: "John", TeamName: "John FC"}}

func TestManagerServiceUpdateInfos(t *testing.T) {
	// TODO add test
}

func TestManagerServiceUpdateTeams(t *testing.T) {
	// TODO add test
}

func TestManagerServiceUpdatePoints(t *testing.T) {
	// TODO add test
}

func TestManagerServiceGetByID(t *testing.T) {
	testcases := []struct {
		ID   int
		want domain.Manager
		err  error
	}{
		{
			ID:   0,
			want: domain.Manager{},
			err:  ErrManagerIDInvalid,
		},
		{
			ID:   john.ID,
			want: john,
			err:  nil,
		},
		{
			ID:   2,
			want: domain.Manager{},
			err:  storage.ErrManagerNotFound,
		},
	}

	mr := mock.ManagerRepository{
		GetByIDFn: func(id int) (domain.Manager, error) {
			if id == john.ID {
				return john, nil
			} else if id == 2 {
				return domain.Manager{}, storage.ErrManagerNotFound
			}

			return domain.Manager{}, nil
		},
	}
	ms := managerService{
		mr: &mr,
	}

	for _, test := range testcases {
		got, err := ms.GetByID(test.ID)
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestManagerServiceConvertToDomainManager(t *testing.T) {
	// TODO add test
}

func TestManagerServiceConvertToDomainTeam(t *testing.T) {
	// TODO add test
}

func TestUpdateTeamPlayersStats(t *testing.T) {
	// TODO add test
}
