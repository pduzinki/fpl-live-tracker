package club

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/mock"
	"testing"
)

var che = domain.Club{ID: 6, Name: "Chelsea", Shortname: "CHE"}

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
			want:   che,
			err:    nil,
		},
		{
			clubID: 66,
			want:   domain.Club{},
			err:    ErrClubIDInvalid,
		},
	}

	cr := mock.ClubRepository{
		GetByIDFn: func(id int) (domain.Club, error) {
			if id != 6 {
				t.Fatalf("unexpected club id: %d", id)
			}
			return che, nil
		},
	}
	cs := clubService{&cr}

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
