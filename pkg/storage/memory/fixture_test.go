package memory

import (
	domain "fpl-live-tracker/pkg"
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
		{fixture: mcitot, want: ErrFixtureAlreadyExists},
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
				t.Errorf("error: incorrect fixture in memory storage")
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
		{},
		{},
	}

	fr := fixtureRepository{
		fixtures: map[int]domain.Fixture{
			mcitot.ID: mcitot,
		},
	}

	for _, test := range testcases {
		got := fr.AddMany(test.fixtures)
		if got != test.want {
			t.Errorf("error: 1")
		}
	}

}

func TestFixtureGetByGameweek(t *testing.T) {
	// TODO
}
