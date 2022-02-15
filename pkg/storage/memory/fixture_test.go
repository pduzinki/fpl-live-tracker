package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
	"time"
)

var (
	livche = domain.Fixture{ID: 1, Info: domain.FixtureInfo{GameweekID: 12, KickoffTime: time.Date(2022, 02, 02, 15, 00, 00, 00, time.UTC)}}
	mcitot = domain.Fixture{ID: 2, Info: domain.FixtureInfo{GameweekID: 13, KickoffTime: time.Date(2022, 02, 02, 16, 00, 00, 00, time.UTC)}}
	burlei = domain.Fixture{ID: 3, Info: domain.FixtureInfo{GameweekID: 13, KickoffTime: time.Date(2022, 02, 03, 15, 00, 00, 00, time.UTC)}}
	cheliv = domain.Fixture{ID: 4, Info: domain.FixtureInfo{GameweekID: 14, KickoffTime: time.Date(2022, 02, 03, 20, 00, 00, 00, time.UTC)}}
)

func TestFixtureAdd(t *testing.T) {
	testcases := []struct {
		fixture domain.Fixture
		want    error
	}{
		{fixture: livche, want: nil},
		{fixture: mcitot, want: storage.ErrFixtureAlreadyExists},
	}

	fr := fixtureRepository{
		fixtures: map[int]domain.Fixture{
			mcitot.ID: mcitot,
		},
	}

	for _, test := range testcases {
		got := fr.Add(test.fixture)
		if got != test.want {
			t.Errorf("error: for %v, got err '%v', want '%v", test.fixture, got, test.want)
		}

		if v, ok := fr.fixtures[test.fixture.ID]; ok {
			if !reflect.DeepEqual(v, test.fixture) {
				t.Errorf("error: incorrect fixture data in memory storage")
			}
		} else {
			t.Errorf("error: fixture not found in memory storage")
		}
	}
}

func TestFixtureAddMany(t *testing.T) {
	testcases := []struct {
		fixtures []domain.Fixture
		want     error
	}{
		{[]domain.Fixture{livche, cheliv}, nil},
		{[]domain.Fixture{cheliv, cheliv}, storage.ErrFixtureAlreadyExists},
	}

	fr := fixtureRepository{
		fixtures: map[int]domain.Fixture{
			mcitot.ID: mcitot,
		},
	}

	for _, test := range testcases {
		got := fr.AddMany(test.fixtures)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}

		for _, f := range test.fixtures {
			if v, ok := fr.fixtures[f.ID]; ok {
				if !reflect.DeepEqual(v, f) {
					t.Errorf("error: incorrect fixture data in memory storage")
				}
			} else {
				t.Errorf("error: fixture not found in memory storage")
			}
		}
	}
}

func TestFixtureUpdate(t *testing.T) {
	testcases := []struct {
		fixture domain.Fixture
		want    error
	}{
		{
			fixture: domain.Fixture{
				ID: mcitot.ID,
				Info: domain.FixtureInfo{
					GameweekID:  22,
					KickoffTime: mcitot.Info.KickoffTime,
				},
			},
			want: nil,
		},
		{
			fixture: domain.Fixture{},
			want:    storage.ErrFixtureNotFound,
		},
	}

	fr := fixtureRepository{
		fixtures: map[int]domain.Fixture{
			livche.ID: livche,
			mcitot.ID: mcitot,
			burlei.ID: burlei,
			cheliv.ID: cheliv,
		},
	}

	for _, test := range testcases {
		got := fr.Update(test.fixture)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}

		if got == nil {
			if v, ok := fr.fixtures[test.fixture.ID]; ok {
				if !reflect.DeepEqual(v, test.fixture) {
					t.Errorf("error: incorrect fixture data in memory storage")
				}
			} else {
				t.Errorf("error: fixture not found in memory storage")
			}
		}
	}
}

func TestFixtureGetByGameweek(t *testing.T) {
	testcases := []struct {
		gw   int
		want []domain.Fixture
		err  error
	}{
		{13, []domain.Fixture{mcitot, burlei}, nil},
		{12, []domain.Fixture{livche}, nil},
		{40, nil, storage.ErrFixtureNotFound},
	}

	fr := fixtureRepository{
		fixtures: map[int]domain.Fixture{
			livche.ID: livche,
			mcitot.ID: mcitot,
			burlei.ID: burlei,
			cheliv.ID: cheliv,
		},
	}

	for _, test := range testcases {
		got, err := fr.GetByGameweek(test.gw)
		if err != test.err {
			t.Errorf("error: got err '%v', want '%v'", err, test.err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: for gw %v, got fixtures %v, want %v", test.gw, got, test.want)
		}
	}
}

func TestFixtureGetByID(t *testing.T) {
	testcases := []struct {
		ID   int
		want domain.Fixture
		err  error
	}{
		{
			ID:   livche.ID,
			want: livche,
			err:  nil,
		},
		{
			ID:   mcitot.ID,
			want: domain.Fixture{},
			err:  storage.ErrFixtureNotFound,
		},
	}

	fr := fixtureRepository{
		fixtures: map[int]domain.Fixture{
			livche.ID: livche,
		},
	}

	for _, test := range testcases {
		got, err := fr.GetByID(test.ID)
		if err != test.err {
			t.Errorf("error: for %v, got '%v', want '%v'", test.ID, got, test.want)
		}
		if err != test.err {
			t.Errorf("error: for %v, got err '%v', want err '%v'", test.ID, err, test.err)
		}
	}
}
