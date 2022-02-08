package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
	"time"
)

var (
	livche = domain.Fixture{Info: domain.FixtureInfo{GameweekID: 12, ID: 1, KickoffTime: time.Date(2022, 02, 02, 15, 00, 00, 00, time.UTC)}}
	mcitot = domain.Fixture{Info: domain.FixtureInfo{GameweekID: 13, ID: 2, KickoffTime: time.Date(2022, 02, 02, 16, 00, 00, 00, time.UTC)}}
	burlei = domain.Fixture{Info: domain.FixtureInfo{GameweekID: 13, ID: 3, KickoffTime: time.Date(2022, 02, 03, 15, 00, 00, 00, time.UTC)}}
	cheliv = domain.Fixture{Info: domain.FixtureInfo{GameweekID: 14, ID: 4, KickoffTime: time.Date(2022, 02, 03, 20, 00, 00, 00, time.UTC)}}
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
			mcitot.Info.ID: mcitot,
		},
	}

	for _, test := range testcases {
		got := fr.Add(test.fixture)
		if got != test.want {
			t.Errorf("error: for %v, got err '%v', want '%v", test.fixture, got, test.want)
		}

		if v, ok := fr.fixtures[test.fixture.Info.ID]; ok {
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
			mcitot.Info.ID: mcitot,
		},
	}

	for _, test := range testcases {
		got := fr.AddMany(test.fixtures)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}

		for _, f := range test.fixtures {
			if v, ok := fr.fixtures[f.Info.ID]; ok {
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
				Info: domain.FixtureInfo{
					GameweekID:  22,
					ID:          mcitot.Info.ID,
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
			livche.Info.ID: livche,
			mcitot.Info.ID: mcitot,
			burlei.Info.ID: burlei,
			cheliv.Info.ID: cheliv,
		},
	}

	for _, test := range testcases {
		got := fr.Update(test.fixture)
		if got != test.want {
			t.Errorf("error: got err '%v', want '%v'", got, test.want)
		}

		if got == nil {
			if v, ok := fr.fixtures[test.fixture.Info.ID]; ok {
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
			livche.Info.ID: livche,
			mcitot.Info.ID: mcitot,
			burlei.Info.ID: burlei,
			cheliv.Info.ID: cheliv,
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
	// TODO add test
}
