package memory

import (
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/storage"
	"testing"
)

var (
	ars = domain.Club{ID: 1, Name: "Arsenal", Shortname: "ARS"}
	che = domain.Club{ID: 2, Name: "Chelsea", Shortname: "CHE"}
	liv = domain.Club{ID: 3, Name: "Liverpool", Shortname: "LIV"}
	mci = domain.Club{ID: 4, Name: "Manchester City", Shortname: "MCI"}
)

func TestClubAdd(t *testing.T) {
	testcases := []struct {
		club domain.Club
		err  error
	}{
		{liv, nil},
		{ars, storage.ErrClubAlreadyExists},
	}

	cr := clubRepository{
		clubs: map[int]domain.Club{
			ars.ID: ars,
		},
	}

	for _, test := range testcases {
		err := cr.Add(test.club)
		if err != test.err {
			t.Errorf("error: for %v, want %v, got %v", test.club, test.err, err)
		}

		if v, ok := cr.clubs[test.club.ID]; ok {
			if v != test.club {
				t.Errorf("error: incorrect club data in memory storage")
			}
		} else {
			t.Errorf("error: club not found in memory storage")
		}
	}
}

func TestClubAddMany(t *testing.T) {
	testcases := []struct {
		clubs []domain.Club
		err   error
	}{
		{[]domain.Club{ars}, nil},
		{[]domain.Club{liv, che}, nil},
		{[]domain.Club{liv}, storage.ErrClubAlreadyExists},
	}

	cr := clubRepository{
		clubs: make(map[int]domain.Club),
	}

	for _, test := range testcases {
		err := cr.AddMany(test.clubs)
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}

		for _, c := range test.clubs {
			if v, ok := cr.clubs[c.ID]; ok {
				if v != c {
					t.Errorf("error: incorrect club data in memory storage")
				}
			} else {
				t.Errorf("error: club not found in memory storage")
			}
		}
	}
}

func TestClubGetByID(t *testing.T) {
	testcases := []struct {
		id   int
		want domain.Club
		err  error
	}{
		{ars.ID, ars, nil},
		{123, domain.Club{}, storage.ErrClubNotFound},
	}

	cr := clubRepository{
		clubs: map[int]domain.Club{
			ars.ID: ars,
		},
	}

	for _, test := range testcases {
		got, err := cr.GetByID(test.id)
		if err != test.err {
			t.Errorf("error: for id %d, want err %v, got %v", test.id, test.err, err)
		}

		if got != test.want {
			t.Errorf("error: for id %d, want %v, got %v", test.id, test.want, got)
		}
	}
}
