package player

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"testing"
)

var ronaldo = domain.Player{
	ID:       123,
	Name:     "Ronaldo",
	Position: "FWD",
}

func TestUpdate(t *testing.T) {
	// TODO add test
}

func TestUpdateStats(t *testing.T) {
	// TODO add test
}

func TestGetByID(t *testing.T) {
	testcases := []struct {
		playerID int
		want     domain.Player
		err      error
	}{
		{
			playerID: 0,
			want:     domain.Player{},
			err:      ErrPlayerIDInvalid,
		},
		{
			playerID: 123,
			want:     ronaldo,
			err:      nil,
		},
	}

	pr := mock.PlayerRepository{
		GetByIDFn: func(id int) (domain.Player, error) {
			if id != 123 {
				t.Fatalf("unexpected id")
			}
			return ronaldo, nil
		},
	}

	ps := playerService{pr: &pr}

	for _, test := range testcases {
		got, err := ps.GetByID(test.playerID)
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}
		if got != test.want {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestConvertToDomainPlayer(t *testing.T) {
	// TODO add test
}

func TestConvertToDomainPlayerStats(t *testing.T) {
	// TODO add test
}
