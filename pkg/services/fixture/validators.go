package fixture

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
)

var ErrGameweekIDInvalid = errors.New("invalid gameweek ID")

type fixtureValidatorFunc func(*domain.Fixture) error

func runFixtureValidations(fixture *domain.Fixture, fns ...fixtureValidatorFunc) error {
	for _, fn := range fns {
		err := fn(fixture)
		if err != nil {
			return err
		}
	}
	return nil
}

func gameweekIDBetween1and38(fixture *domain.Fixture) error {
	if fixture.Info.GameweekID <= 0 || fixture.Info.GameweekID > 38 {
		return ErrGameweekIDInvalid
	}

	return nil
}
