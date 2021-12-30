package memory

import (
	"errors"
	domain "fpl-live-tracker/pkg"
	"sync"
)

var (
	ErrFixtureNotFound      error = errors.New("storage: fixture not found")
	ErrFixtureAlreadyExists error = errors.New("storage: fixture already exists")
)

//
type fixtureRepository struct {
	fixtures map[int]domain.Fixture // TODO change to more proper data structure
	sync.Mutex
}

//
func NewFixtureRepository() domain.FixtureRepository {
	return &fixtureRepository{
		fixtures: make(map[int]domain.Fixture),
	}
}

//
func (fr *fixtureRepository) Add(fixture domain.Fixture) error {
	if _, ok := fr.fixtures[fixture.ID]; ok {
		return ErrFixtureAlreadyExists
	}

	fr.Lock()
	fr.fixtures[fixture.ID] = fixture
	fr.Unlock()

	return nil
}

//
func (fr *fixtureRepository) AddMany(fixtures []domain.Fixture) error {
	for _, fixture := range fixtures {
		err := fr.Add(fixture)
		if err != nil {
			return err
		}
	}
	return nil
}

//
func (fr *fixtureRepository) GetByGameweek(gameweekID int) ([]domain.Fixture, error) {
	fixtures := make([]domain.Fixture, 0)

	for _, f := range fr.fixtures {
		if f.GameweekID == gameweekID {
			fixtures = append(fixtures, f)
		}
	}

	return fixtures, nil
}
