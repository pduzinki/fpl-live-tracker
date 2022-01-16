package club

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
)

var ErrClubIDInvalid = errors.New("invalid club ID")

type clubValidatorFunc func(*domain.Club) error

func runClubValidations(club *domain.Club, fns ...clubValidatorFunc) error {
	for _, fn := range fns {
		err := fn(club)
		if err != nil {
			return err
		}
	}
	return nil
}

func idBetween1and20(club *domain.Club) error {
	if club.ID <= 0 || club.ID > 20 {
		return ErrClubIDInvalid
	}

	return nil
}
