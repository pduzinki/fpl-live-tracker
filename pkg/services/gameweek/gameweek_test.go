package gameweek

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"fpl-live-tracker/pkg/wrapper"
	"testing"
)

var (
	errWrapperFail = errors.New("wrapper fail")
)

var gameweek = domain.Gameweek{
	ID:       4,
	Name:     "Gameweek 4",
	Finished: false,
}

func TestUpdate(t *testing.T) {
	testcases := []struct {
		name string
		err  error
	}{
		{
			name: "wrapper fail",
			err:  errWrapperFail,
		},
		{
			name: "sunny scenario",
			err:  nil,
		},
	}

	for _, test := range testcases {
		wr := mock.Wrapper{
			GetGameweeksFn: func() ([]wrapper.Gameweek, error) {
				if test.name == "wrapper fail" {
					return []wrapper.Gameweek{}, errWrapperFail
				} else if test.name == "current gw conversion fail" {
					return []wrapper.Gameweek{
						{
							ID:           4,
							Name:         "Gameweek 4",
							Finished:     false,
							IsCurrent:    true,
							IsNext:       false,
							DeadlineTime: "broken time",
						},
					}, nil
				} else if test.name == "sunny scenario" {
					return []wrapper.Gameweek{
						{
							ID:           4,
							Name:         "Gameweek 4",
							Finished:     false,
							IsCurrent:    true,
							IsNext:       false,
							DeadlineTime: "2021-12-04T12:30:00Z",
						},
						{
							ID:           5,
							Name:         "Gameweek 5",
							Finished:     false,
							IsCurrent:    false,
							IsNext:       true,
							DeadlineTime: "2021-12-12T12:30:00Z",
						},
					}, nil
				}
				t.Fatalf("unexpected test name")
				return nil, nil
			},
		}

		gs := gameweekService{wr: &wr}

		err := gs.Update()
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)

		}
	}
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

func TestConvertToDomainGameweek(t *testing.T) {
	// TODO add test
}
