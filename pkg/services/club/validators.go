package club

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
)

var ErrClubIDInvalid = errors.New("invalid id")

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

func idGreaterThanZero(club *domain.Club) error {
	if club.ID <= 0 {
		return ErrClubIDInvalid
	}

	return nil
}

func idNotGreaterThanTwenty(club *domain.Club) error {
	if club.ID > 20 {
		return ErrClubIDInvalid
	}

	return nil
}
