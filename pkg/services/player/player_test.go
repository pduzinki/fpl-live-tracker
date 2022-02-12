package player

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"reflect"
	"testing"
)

var ronaldo = domain.Player{
	ID:       123,
	Name:     "Ronaldo",
	Position: "FWD",
}

func TestUpdate(t *testing.T) {
	// TODO add test
}

func TestUpdateStats(t *testing.T) {
	// TODO add test
}

func TestGetByID(t *testing.T) {
	testcases := []struct {
		playerID int
		want     domain.Player
		err      error
	}{
		{
			playerID: 0,
			want:     domain.Player{},
			err:      ErrPlayerIDInvalid,
		},
		{
			playerID: 123,
			want:     ronaldo,
			err:      nil,
		},
	}

	pr := mock.PlayerRepository{
		GetByIDFn: func(id int) (domain.Player, error) {
			if id != 123 {
				t.Fatalf("unexpected id")
			}
			return ronaldo, nil
		},
	}

	ps := playerService{pr: &pr}

	for _, test := range testcases {
		got, err := ps.GetByID(test.playerID)
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestConvertToDomainPlayer(t *testing.T) {
	// TODO add test
}

func TestConvertToDomainPlayerStats(t *testing.T) {
	// TODO add test
}

func TestFindTopBPS(t *testing.T) {
	testcases := []struct {
		data []domain.FixtureStatPair
		want []int
	}{
		{
			data: []domain.FixtureStatPair{{PlayerID: 1, Value: 55}, {PlayerID: 2, Value: 44},
				{PlayerID: 3, Value: 33}, {PlayerID: 4, Value: 33},
				{PlayerID: 5, Value: 33}, {PlayerID: 6, Value: 22}},
			want: []int{55, 44, 33},
		},
		{
			data: []domain.FixtureStatPair{{PlayerID: 1, Value: 55}, {PlayerID: 2, Value: 55},
				{PlayerID: 3, Value: 55}, {PlayerID: 4, Value: 55},
				{PlayerID: 5, Value: 33}, {PlayerID: 6, Value: 22}},
			want: []int{55, 33, 22},
		},
		{
			data: []domain.FixtureStatPair{{PlayerID: 1, Value: 5}, {PlayerID: 2, Value: 5},
				{PlayerID: 3, Value: 5}, {PlayerID: 4, Value: 5},
				{PlayerID: 5, Value: 5}, {PlayerID: 6, Value: 5}},
			want: []int{5},
		},
		{
			data: []domain.FixtureStatPair{},
			want: []int{},
		},
	}

	for _, test := range testcases {
		got := findTopBPS(test.data)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestFindPlayersAndBonusPoints(t *testing.T) {
	// TODO add test
}

func TestAddBonusPoints(t *testing.T) {
	// TODO add test
}
