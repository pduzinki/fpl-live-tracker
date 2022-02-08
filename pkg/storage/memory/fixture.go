package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sort"
	"sync"
)

//
type fixtureRepository struct {
	fixtures map[int]domain.Fixture
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
	if _, ok := fr.fixtures[fixture.Info.ID]; ok {
		return storage.ErrFixtureAlreadyExists
	}

	fr.Lock()
	fr.fixtures[fixture.Info.ID] = fixture
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
func (fr *fixtureRepository) Update(fixture domain.Fixture) error {
	if _, ok := fr.fixtures[fixture.Info.ID]; ok {
		fr.fixtures[fixture.Info.ID] = fixture
		return nil
	}

	return storage.ErrFixtureNotFound
}

//
func (fr *fixtureRepository) GetByGameweek(gameweekID int) ([]domain.Fixture, error) {
	fixtures := make([]domain.Fixture, 0)

	for _, f := range fr.fixtures {
		if f.Info.GameweekID == gameweekID {
			fixtures = append(fixtures, f)
		}
	}

	if len(fixtures) == 0 {
		return nil, storage.ErrFixtureNotFound
	}

	sort.Slice(fixtures, func(i, j int) bool {
		return fixtures[i].Info.KickoffTime.Before(fixtures[j].Info.KickoffTime)
	})

	return fixtures, nil
}

//
func (fr *fixtureRepository) GetByID(fixtureID int) (domain.Fixture, error) {
	if fixture, ok := fr.fixtures[fixtureID]; ok {
		return fixture, nil
	}
	return domain.Fixture{}, storage.ErrFixtureNotFound
}
