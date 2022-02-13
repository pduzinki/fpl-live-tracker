package manager

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
)

var ErrManagerIDInvalid = errors.New("invalid manager ID")

type managerValidatorFunc func(*domain.Manager) error

func runManagerValidations(manager *domain.Manager, fns ...managerValidatorFunc) error {
	for _, fn := range fns {
		err := fn(manager)
		if err != nil {
			return err
		}
	}
	return nil
}

func idHigherThanZero(manager *domain.Manager) error {
	if manager.ID <= 0 {
		return ErrManagerIDInvalid
	}
	return nil
}
