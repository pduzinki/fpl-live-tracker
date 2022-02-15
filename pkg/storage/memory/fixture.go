package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sort"
	"sync"
)

// fixtureRepository implements domain.FixtureRepository interface
type fixtureRepository struct {
	fixtures map[int]domain.Fixture
	sync.RWMutex
}

// NewFixtureRepository returns new instance of domain.FixtureRepository
func NewFixtureRepository() domain.FixtureRepository {
	return &fixtureRepository{
		fixtures: make(map[int]domain.Fixture),
	}
}

// Add saves given fixture into memory storage, returns error on failure
func (fr *fixtureRepository) Add(fixture domain.Fixture) error {
	fr.Lock()
	defer fr.Unlock()

	if _, ok := fr.fixtures[fixture.Info.ID]; ok {
		return storage.ErrFixtureAlreadyExists
	}
	fr.fixtures[fixture.Info.ID] = fixture

	return nil
}

// AddMany saves given fixtures into memory storage, returns error on failure
func (fr *fixtureRepository) AddMany(fixtures []domain.Fixture) error {
	for _, fixture := range fixtures {
		err := fr.Add(fixture)
		if err != nil {
			return err
		}
	}
	return nil
}

// Update updates fixture with matching ID in memory storage, or returns error on failure
func (fr *fixtureRepository) Update(fixture domain.Fixture) error {
	fr.Lock()
	defer fr.Unlock()

	if _, ok := fr.fixtures[fixture.Info.ID]; ok {
		fr.fixtures[fixture.Info.ID] = fixture
		return nil
	}

	return storage.ErrFixtureNotFound
}

// GetByGameweek returns all fixtures with given gameweekID, or error otherwise.
// Returned fixtures will be sorted by kick-off time
func (fr *fixtureRepository) GetByGameweek(gameweekID int) ([]domain.Fixture, error) {
	fr.RLock()

	fixtures := make([]domain.Fixture, 0)

	for _, f := range fr.fixtures {
		if f.Info.GameweekID == gameweekID {
			fixtures = append(fixtures, f)
		}
	}
	fr.RUnlock()

	if len(fixtures) == 0 {
		return nil, storage.ErrFixtureNotFound
	}

	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].Info.KickoffTime.Before(fixtures[j].Info.KickoffTime)
	})

	return fixtures, nil
}

// GetById returns fixtures with given ID, or returns error otherwise
func (fr *fixtureRepository) GetByID(fixtureID int) (domain.Fixture, error) {
	fr.RLock()
	defer fr.RUnlock()

	if fixture, ok := fr.fixtures[fixtureID]; ok {
		return fixture, nil
	}
	return domain.Fixture{}, storage.ErrFixtureNotFound
}
