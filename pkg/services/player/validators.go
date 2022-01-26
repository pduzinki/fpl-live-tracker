package player

import (
	"errors"
	"fpl-live-tracker/pkg/domain"
)

var ErrPlayerIDInvalid = errors.New("invalid player ID")

type playerValidatorFunc func(*domain.Player) error

func runPlayerValidations(player *domain.Player, fns ...playerValidatorFunc) error {
	for _, fn := range fns {
		err := fn(player)
		if err != nil {
			return err
		}
	}
	return nil
}

func idHigherThanZero(player *domain.Player) error {
	if player.ID <= 0 {
		return ErrPlayerIDInvalid
	}
	return nil
}
