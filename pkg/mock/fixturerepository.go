package mock

import "fpl-live-tracker/pkg/domain"

type FixtureRepository struct {
	AddFn           func(fixture domain.Fixture) error
	AddManyFn       func(fixtures []domain.Fixture) error
	UpdateFn        func(fixture domain.Fixture) error
	GetByGameweekFn func(gameweekID int) ([]domain.Fixture, error)
	GetByIDFn       func(fixtureID int) (domain.Fixture, error)
}

func (fr *FixtureRepository) Add(fixture domain.Fixture) error {
	return fr.AddFn(fixture)
}

func (fr *FixtureRepository) AddMany(fixtures []domain.Fixture) error {
	return fr.AddManyFn(fixtures)
}

func (fr *FixtureRepository) Update(fixture domain.Fixture) error {
	return fr.UpdateFn(fixture)
}

func (fr *FixtureRepository) GetByGameweek(gameweekID int) ([]domain.Fixture, error) {
	return fr.GetByGameweekFn(gameweekID)
}

func (fr *FixtureRepository) GetByID(fixtureID int) (domain.Fixture, error) {
	return fr.GetByIDFn(fixtureID)
}
