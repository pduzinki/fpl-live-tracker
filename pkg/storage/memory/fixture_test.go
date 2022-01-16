package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"
)

var (
	livche = domain.Fixture{GameweekID: 12, ID: 1}
	mcitot = domain.Fixture{GameweekID: 13, ID: 2}
	burlei = domain.Fixture{GameweekID: 13, ID: 3}
	cheliv = domain.Fixture{GameweekID: 14, ID: 4}
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
			if v != test.fixture {
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
				if v != f {
					t.Errorf("error: incorrect fixture data in memory storage")
				}
			} else {
				t.Errorf("error: fixture not found in memory storage")
			}
		}
	}
}

func TestFixtureUpdate(t *testing.T) {
	// TODO add test
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
