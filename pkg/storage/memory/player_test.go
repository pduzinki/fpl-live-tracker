package memory

import (
	"fpl-live-tracker/pkg/storage"
	"reflect"
	"testing"

	domain "fpl-live-tracker/pkg"
)

var (
	ramsdale = domain.Player{
		ID:       1,
		Name:     "Ramsdale",
		Position: "GKP",
		Club:     domain.Club{ID: 1, Name: "Arsenal", Shortname: "ARS"},
	}
	cancelo = domain.Player{
		ID:       2,
		Name:     "Cancelo",
		Position: "DEF",
		Club:     domain.Club{ID: 2, Name: "Manchester City", Shortname: "MCI"},
	}
	salah = domain.Player{
		ID:       3,
		Name:     "Salah",
		Position: "MID",
		Club:     domain.Club{ID: 3, Name: "Liverpool", Shortname: "LIV"},
	}
	kane = domain.Player{
		ID:       4,
		Name:     "Kane",
		Position: "FWD",
		Club:     domain.Club{ID: 4, Name: "Spurs", Shortname: "TOT"},
	}
)

func TestPlayerAdd(t *testing.T) {
	testcases := []struct {
		player domain.Player
		want   error
	}{
		{player: salah, want: nil},
		{player: kane, want: storage.ErrPlayerAlreadyExists},
	}

	pr := playerRepository{
		players: map[int]domain.Player{
			kane.ID: kane,
		},
	}

	for _, test := range testcases {
		got := pr.Add(test.player)
		if got != test.want {
			t.Errorf("error: for %v, got err '%v', want '%v", test.player, got, test.want)
		}

		if v, ok := pr.players[test.player.ID]; ok {
			if v != test.player {
				t.Errorf("error: incorrect player data in memory storage")
			}
		} else {
			t.Errorf("error: player not found in memory storage")
		}
	}
}

func TestPlayerUpdate(t *testing.T) {
	// TODO add test
}

func TestPlayerGetByID(t *testing.T) {
	testcases := []struct {
		id   int
		want domain.Player
		err  error
	}{
		{kane.ID, kane, nil},
		{81, domain.Player{}, storage.ErrPlayerNotFound},
	}

	pr := playerRepository{
		players: map[int]domain.Player{
			kane.ID:     kane,
			salah.ID:    salah,
			cancelo.ID:  cancelo,
			ramsdale.ID: ramsdale,
		},
	}

	for _, test := range testcases {
		got, err := pr.GetByID(test.id)
		if err != test.err {
			t.Errorf("error: got err '%v', want '%v'", err, test.err)
		}

		if !reflect.DeepEqual(got, test.want) {
			t.Errorf("error: for player %v, got player %v, want %v", test.id, got, test.want)
		}
	}
}
