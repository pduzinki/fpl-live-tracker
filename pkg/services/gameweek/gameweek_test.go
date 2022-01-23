package gameweek

import (
	"fpl-live-tracker/pkg/domain"
	"testing"
)

var gameweek = domain.Gameweek{
	ID:       4,
	Name:     "Gameweek 4",
	Finished: false,
}

func TestUpdate(t *testing.T) {
	// TODO add test
}

func TestGetCurrentGameweek(t *testing.T) {
	testcases := []struct {
		want    domain.Gameweek
		err     error
		updated bool
	}{
		{
			want:    gameweek,
			err:     nil,
			updated: true,
		},
		{
			want:    domain.Gameweek{},
			err:     ErrGameweekNotUpdated,
			updated: false,
		},
	}

	for _, test := range testcases {
		gs := gameweekService{
			CurrentGameweek: gameweek,
			updatedOnce:     test.updated,
		}

		got, err := gs.GetCurrentGameweek()
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}
		if got != test.want {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestGetNextGameweek(t *testing.T) {
	testcases := []struct {
		want         domain.Gameweek
		err          error
		updated      bool
		gameFinished bool
	}{
		{
			want:         gameweek,
			err:          nil,
			updated:      true,
			gameFinished: false,
		},
		{
			want:         domain.Gameweek{},
			err:          ErrGameweekNotUpdated,
			updated:      false,
			gameFinished: false,
		},
		{
			want:         domain.Gameweek{},
			err:          ErrNoNextGameweek,
			updated:      true,
			gameFinished: true,
		},
	}

	for _, test := range testcases {
		gs := gameweekService{
			NextGameweek:   gameweek,
			updatedOnce:    test.updated,
			noNextGameweek: test.gameFinished,
		}

		got, err := gs.GetNextGameweek()
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}
		if got != test.want {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}
