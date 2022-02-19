package manager

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
)

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

func TestCalculateTotalPoints(t *testing.T) {
	testcases := []struct {
		name string
		team domain.Team
		want int
	}{
		{
			name: "empty team passed",
			team: domain.Team{},
			want: 0,
		},
		{
			name: "no chip played",
			team: teamA,
			want: 66,
		},
		{
			name: "triple captain played",
			team: chipTripleCpt(teamA),
			want: 81,
		},
		{
			name: "bench boost played",
			team: chipBenchBoost(teamA),
			want: 84,
		},
	}

	for _, test := range testcases {
		got := calculateTotalPoints(&test.team)
		if got != test.want {
			t.Errorf("error: for test '%s': want %v, got %v", test.name, test.want, got)
		}
	}
}

func TestCalculateSubPoints(t *testing.T) {
	testcases := []struct {
		name string
		team domain.Team
		want int
	}{
		{
			name: "empty team passed",
			team: domain.Team{},
			want: 0,
		},
		{
			name: "no subs",
			team: teamA,
			want: 0,
		},
		{
			name: "two subs",
			team: teamB,
			want: 16,
		},
		{
			name: "all subs needed",
			team: teamD,
			want: 18,
		},
		{
			name: "def subs needed",
			team: teamE,
			want: 2,
		},
		{
			name: "fwd sub needed",
			team: teamF,
			want: 1,
		},
		{
			name: "bench boost played",
			team: chipBenchBoost(teamA),
			want: 0,
		},
		{
			name: "captain sub",
			team: teamC,
			want: 8,
		},
		{
			name: "triple captain sub",
			team: chipTripleCpt(teamC),
			want: 15,
		},
		{
			name: "too few defs, formation not broken",
			team: teamG,
			want: 0,
		},
		{
			name: "too few fwds, formation not broken",
			team: teamH,
			want: 0,
		},
	}

	for _, test := range testcases {
		got := calculateSubPoints(&test.team)
		if got != test.want {
			t.Errorf("error: for test '%s': want %v, got %v", test.name, test.want, got)
		}
	}
}

func TestGetLiveFormation(t *testing.T) {
	// TODO add test
}

func TestPlayed(t *testing.T) {
	// TODO add test
}

func TestCaptainPlayed(t *testing.T) {
	// TODO add test
}
