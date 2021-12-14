package fixture

import (
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/wrapper"
)

type FixtureService interface {
	GetFixtures(gameweekID int) ([]domain.Fixture, error)
}

type fixtureService struct {
	// wrapper wrapper.Wrapper
}

//
func NewFixtureService(w wrapper.Wrapper) FixtureService {
	return &fixtureService{}
}

//
func (fs *fixtureService) GetFixtures(gameweekID int) ([]domain.Fixture, error) {
	return nil, nil
}
