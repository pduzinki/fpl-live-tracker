package club

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"fpl-live-tracker/pkg/storage/memory"
	"testing"
)

func TestGetClubByID(t *testing.T) {
	testcases := []struct {
		clubID int
		want   domain.Club
		err    error
	}{
		{
			clubID: 0,
			want:   domain.Club{},
			err:    ErrClubIDInvalid,
		},
		{
			clubID: 6,
			want:   domain.Club{ID: 6, Name: "Chelsea", Shortname: "CHE"},
			err:    nil,
		},
		{
			clubID: 66,
			want:   domain.Club{},
			err:    ErrClubIDInvalid,
		},
	}

	cr := memory.NewClubRepository()

	wr := mock.Wrapper{
		GetClubsFn: mock.GetClubsOK,
	}

	cs, err := NewClubService(cr, &wr)
	if err != nil {
		t.Fatal("error: failed to init club service")
	}

	for _, test := range testcases {
		got, err := cs.GetClubByID(test.clubID)
		if err != test.err {
			t.Errorf("error: want err %v, got %v", test.err, err)
		}

		if got != test.want {
			t.Errorf("error: want %v, got %v", test.want, got)
		}
	}
}
