package player

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
	errWrapperFail        = errors.New("wrapper fail")
	errGetClubFail        = errors.New("get club fail")
	errRepoUpdateInfoFail = errors.New("update info fail")
	errAddFail            = errors.New("add player fail")
)

var wrapperPlayers = []wrapper.Player{
	{
		ID:       1,
		Team:     1,
		Position: 1,
		WebName:  "Ramsdale",
	},
	{
		ID:       2,
		Team:     2,
		Position: 3,
		WebName:  "Coutinho",
	},
}

var ronaldo = domain.Player{
	ID: 123,
	Info: domain.PlayerInfo{
		Name:     "Ronaldo",
		Position: "FWD",
	},
}

func TestUpdateInfos(t *testing.T) {
	testcases := []struct {
		name string
		err  error
	}{
		{
			name: "wrapper fail",
			err:  errWrapperFail,
		},
		{
			name: "get club fail",
			err:  errGetClubFail,
		},
		{
			name: "update info fail",
			err:  errRepoUpdateInfoFail,
		},
		{
			name: "update info player not found fail",
			err:  errAddFail,
		},
		{
			name: "sunny scenario",
			err:  nil,
		},
	}

	for _, test := range testcases {
		wr := mock.Wrapper{
			GetPlayersFn: func() ([]wrapper.Player, error) {
				if test.name == "wrapper fail" {
					return []wrapper.Player{}, errWrapperFail
				}
				return wrapperPlayers, nil
			},
		}

		cs := mock.ClubService{
			GetClubByIDFn: func(id int) (domain.Club, error) {
				if test.name == "get club fail" {
					return domain.Club{}, errGetClubFail
				}
				return domain.Club{}, nil
			},
		}

		pr := mock.PlayerRepository{
			AddFn: func(player domain.Player) error {
				if test.name == "update info player not found fail" {
					return errAddFail
				}
				return nil
			},

			UpdateInfoFn: func(playerID int, info domain.PlayerInfo) error {
				if test.name == "update info fail" {
					return errRepoUpdateInfoFail
				} else if test.name == "update info player not found fail" {
					return storage.ErrPlayerNotFound
				}
				return nil
			},
		}

		ps := playerService{
			wr: &wr,
			cs: &cs,
			pr: &pr,
		}

		err := ps.UpdateInfos()
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}
	}
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
	testcases := []struct {
		stats []domain.FixtureStatPair
		bps   []int
		want  []bonusPlayer
	}{
		{
			stats: []domain.FixtureStatPair{
				{PlayerID: 1, Value: 44},
				{PlayerID: 2, Value: 33},
				{PlayerID: 3, Value: 22}},
			bps:  []int{},
			want: []bonusPlayer{},
		},
		{
			stats: []domain.FixtureStatPair{},
			bps:   []int{44, 33, 22},
			want:  []bonusPlayer{},
		},
		{
			stats: []domain.FixtureStatPair{
				{PlayerID: 1, Value: 44},
				{PlayerID: 2, Value: 33},
				{PlayerID: 3, Value: 22}},
			bps: []int{44, 33, 22},
			want: []bonusPlayer{
				{playerID: 1, bonusPoints: 3},
				{playerID: 2, bonusPoints: 2},
				{playerID: 3, bonusPoints: 1},
			},
		},
		{
			stats: []domain.FixtureStatPair{
				{PlayerID: 1, Value: 44},
				{PlayerID: 2, Value: 44},
				{PlayerID: 3, Value: 44},
				{PlayerID: 4, Value: 33},
				{PlayerID: 5, Value: 22}},
			bps: []int{44, 33, 22},
			want: []bonusPlayer{
				{playerID: 1, bonusPoints: 3},
				{playerID: 2, bonusPoints: 3},
				{playerID: 3, bonusPoints: 3},
			},
		},
		{
			stats: []domain.FixtureStatPair{
				{PlayerID: 1, Value: 44},
				{PlayerID: 2, Value: 33},
				{PlayerID: 3, Value: 33},
				{PlayerID: 4, Value: 33},
				{PlayerID: 5, Value: 22}},
			bps: []int{44, 33, 22},
			want: []bonusPlayer{
				{playerID: 1, bonusPoints: 3},
				{playerID: 2, bonusPoints: 2},
				{playerID: 3, bonusPoints: 2},
				{playerID: 4, bonusPoints: 2},
			},
		},
		{
			stats: []domain.FixtureStatPair{
				{PlayerID: 1, Value: 44},
				{PlayerID: 2, Value: 33},
				{PlayerID: 3, Value: 22},
				{PlayerID: 4, Value: 22}},
			bps: []int{44, 33, 22},
			want: []bonusPlayer{
				{playerID: 1, bonusPoints: 3},
				{playerID: 2, bonusPoints: 2},
				{playerID: 3, bonusPoints: 1},
				{playerID: 4, bonusPoints: 1},
			},
		},
	}

	for _, test := range testcases {
		got := findPlayersAndBonusPoints(test.stats, test.bps)
		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}

func TestAddBonusPoints(t *testing.T) {
	// TODO add test
}
