package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
)

var (
	jimsTeam  = domain.Team{ID: 1}
	janesTeam = domain.Team{ID: 2}
	joelsTeam = domain.Team{ID: 3}
	jacksTeam = domain.Team{ID: 4}
)

func TestTeamAdd(t *testing.T) {
	testcases := []struct {
		team domain.Team
		want error
	}{
		{joelsTeam, nil},
		{jimsTeam, storage.ErrTeamAlreadyExists},
	}

	tr := teamRepository{
		teams: map[int]domain.Team{
			jimsTeam.ID:  jimsTeam,
			janesTeam.ID: janesTeam,
		},
	}

	for _, test := range testcases {
		got := tr.Add(test.team)
		if got != test.want {
			t.Errorf("error: for %v, got err '%v', want '%v'", test.team, got, test.want)
		}

		if v, ok := tr.teams[test.team.ID]; ok {
			if !reflect.DeepEqual(v, test.team) {
				t.Errorf("error: incorrect team data in memory storage")
			}
		} else {
			t.Errorf("error: team not found in memory storage")
		}
	}
}

func TestTeamUpdate(t *testing.T) {
	testcases := []struct {
		teamID int
		team   domain.Team
		want   error
	}{
		{
			teamID: jacksTeam.ID,
			team: domain.Team{
				ID:         jacksTeam.ID,
				GameweekID: 22,
			},
			want: nil,
		},
		{
			teamID: 420,
			team:   domain.Team{},
			want:   storage.ErrTeamNotFound,
		},
	}

	tr := teamRepository{
		teams: map[int]domain.Team{
			jacksTeam.ID: jacksTeam,
		},
	}

	for _, test := range testcases {
		got := tr.Update(test.teamID, test.team)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}

		if got == nil {
			if v, ok := tr.teams[test.teamID]; ok {
				if !reflect.DeepEqual(v, test.team) {
					t.Errorf("error: incorrect team data in memory storage")
				}
			} else {
				t.Errorf("error: team not found in memory storage")
			}
		}

	}
}

func TestTeamGetByID(t *testing.T) {
	// TODO add test
}

func TestTeamGetCount(t *testing.T) {
	// TODO add test
}
