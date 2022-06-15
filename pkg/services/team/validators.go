package team

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
)

var ErrTeamIDInvalid = errors.New("invalid team ID")

type teamValidatorFunc func(*domain.Team) error

func runTeamValidations(team *domain.Team, fns ...teamValidatorFunc) error {
	for _, fn := range fns {
		err := fn(team)
		if err != nil {
			return err
		}
	}
	return nil
}

func idHigherThanZero(team *domain.Team) error {
	if team.ID <= 0 {
		return ErrTeamIDInvalid
	}
	return nil
}
