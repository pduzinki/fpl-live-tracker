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
		want domain.Gameweek
		err  error
	}{
		{
			want: gameweek,
			err:  nil,
		},
	}

	for _, test := range testcases {
		gs := gameweekService{
			CurrentGameweek: gameweek,
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
		gameFinished bool
	}{
		{
			want:         gameweek,
			err:          nil,
			gameFinished: false,
		},
		{
			want:         domain.Gameweek{},
			err:          ErrNoNextGameweek,
			gameFinished: true,
		},
	}

	for _, test := range testcases {
		gs := gameweekService{
			NextGameweek:   gameweek,
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
