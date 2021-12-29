package fixture

import (
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/wrapper"
)

type FixtureService interface {
	GetFixtures(gameweekID int) ([]domain.Fixture, error)
}

type fixtureService struct {
	wrapper wrapper.Wrapper
}

//
func NewFixtureService(w wrapper.Wrapper) FixtureService {
	return &fixtureService{
		wrapper: w,
	}
}

//
func (fs *fixtureService) GetFixtures(gameweekID int) ([]domain.Fixture, error) {
	fixtures, err := fs.wrapper.GetFixtures(gameweekID)
	if err != nil {
		return []domain.Fixture{}, err
	}

	return fixtures, nil
}
