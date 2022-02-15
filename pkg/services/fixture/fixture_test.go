package fixture

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"reflect"
	"testing"
)

var (
	errWrapperFail = errors.New("wrapper fail")
	errGetClubFail = errors.New("get club fail")
	errUpdateFail  = errors.New("update fail")
	errAddFail     = errors.New("add fail")
)

var wrapperFixtures = []wrapper.Fixture{
	{
		Event:               8,
		ID:                  123,
		KickoffTime:         "2021-12-04T12:30:00Z",
		Started:             false,
		Finished:            false,
		FinishedProvisional: false,
		TeamA:               1,
		TeamH:               2,
		Stats:               []wrapper.FixtureStat{},
	},
	{
		Event:               8,
		ID:                  124,
		KickoffTime:         "2021-12-05T12:30:00Z",
		Started:             false,
		Finished:            false,
		FinishedProvisional: false,
		TeamA:               3,
		TeamH:               4,
		Stats:               []wrapper.FixtureStat{},
	},
}

var gw8Fixtures = []domain.Fixture{
	{ID: 123, Info: domain.FixtureInfo{GameweekID: 8, Started: true, FinishedProvisional: true, Finished: true}},
	{ID: 124, Info: domain.FixtureInfo{GameweekID: 8, Started: true, FinishedProvisional: true, Finished: false}},
	{ID: 125, Info: domain.FixtureInfo{GameweekID: 8, Started: true, FinishedProvisional: false, Finished: false}},
	{ID: 126, Info: domain.FixtureInfo{GameweekID: 8, Started: false, FinishedProvisional: false, Finished: false}},
	{ID: 127, Info: domain.FixtureInfo{GameweekID: 8, Started: false, FinishedProvisional: false, Finished: false}},
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
			name: "wrapper type to domain type conversion fail",
			err:  errGetClubFail,
		},
		{
			name: "update fixtures fail",
			err:  errUpdateFail,
		},
		{
			name: "add fixtures fail",
			err:  errAddFail,
		},
		{
			name: "sunny scenario",
			err:  nil,
		},
	}

	for _, test := range testcases {
		fr := mock.FixtureRepository{
			AddFn: func(f domain.Fixture) error {
				if test.name == "add fixtures fail" {
					return errAddFail
				}
				return nil
			},
			UpdateFn: func(f domain.Fixture) error {
				if test.name == "update fixtures fail" {
					return errUpdateFail
				}
				if test.name == "add fixtures fail" {
					return storage.ErrFixtureNotFound
				}
				return nil
			},
		}

		wr := mock.Wrapper{
			GetFixturesFn: func() ([]wrapper.Fixture, error) {
				if test.name == "wrapper fail" {
					return []wrapper.Fixture{}, errWrapperFail
				}
				return wrapperFixtures, nil
			},
		}

		cs := mock.ClubService{
			GetClubByIDFn: func(id int) (domain.Club, error) {
				if test.name == "wrapper type to domain type conversion fail" {
					return domain.Club{}, errGetClubFail
				}
				return domain.Club{}, nil
			},
		}

		fs := fixtureService{&fr, &cs, &wr}

		err := fs.Update()
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}
	}
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
	fs := fixtureService{fr: &fr}

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
			want: gw8Fixtures[1:3],
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
	fs := fixtureService{fr: &fr}

	for _, test := range testcases {
		got, err := fs.GetLiveFixtures(test.gwID)
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestConvertToDomainFixture(t *testing.T) {
	// TODO add test
}
