package fixture

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"reflect"
	"testing"
)

var gw8Fixtures = []domain.Fixture{
	{GameweekID: 8, ID: 123},
	{GameweekID: 8, ID: 124},
}

func TestUpdate(t *testing.T) {
	// TODO add test
}

func TestGetFixturesByGameweek(t *testing.T) {
	testcases := []struct {
		gwID int
		want []domain.Fixture
		err  error
	}{
		{
			gwID: 0,
			want: []domain.Fixture{},
			err:  ErrGameweekIDInvalid,
		},
		{
			gwID: 8,
			want: gw8Fixtures,
			err:  nil,
		},
		{
			gwID: 444,
			want: []domain.Fixture{},
			err:  ErrGameweekIDInvalid,
		},
	}

	fr := mock.FixtureRepository{
		GetByGameweekFn: func(gameweekID int) ([]domain.Fixture, error) {
			return gw8Fixtures, nil
		},
	}
	wr := mock.Wrapper{}
	cs := mock.ClubService{}

	fs := NewFixtureService(&fr, &cs, &wr)

	for _, test := range testcases {
		got, err := fs.GetFixturesByGameweek(test.gwID)
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestGetLiveFixtures(t *testing.T) {
	// TODO add test
}
