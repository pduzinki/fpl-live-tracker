package fixture

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"time"
)

// FixtureService is an interface for interacting with fixtures
type FixtureService interface {
	Update() error
	GetFixturesByGameweek(gameweekID int) ([]domain.Fixture, error)
	GetLiveFixtures(gameweekID int) ([]domain.Fixture, error)
}

// fixtureService implements FixtureService interface
type fixtureService struct {
	fr      domain.FixtureRepository
	cs      club.ClubService
	wrapper wrapper.Wrapper
}

// NewFixtureService returns new instance of FixtureService, and fills
//underlying storage with data from FPL API
func NewFixtureService(fr domain.FixtureRepository, cs club.ClubService,
	w wrapper.Wrapper) (FixtureService, error) {
	fs := fixtureService{
		fr:      fr,
		cs:      cs,
		wrapper: w,
	}

	err := fs.Update()
	if err != nil {
		log.Println("fixture service: failed to init data")
		return nil, err
	}

	return &fs, nil
}

// Update queries FPL API and updates all fixture data in its underlying fixture storage
func (fs *fixtureService) Update() error {
	wrapperFixtures, err := fs.wrapper.GetFixtures()
	if err != nil {
		log.Println("fixture service:", err)
		return err
	}

	fixtures := make([]domain.Fixture, len(wrapperFixtures))
	for i, wf := range wrapperFixtures {
		f, err := fs.convertToDomainFixture(wf)
		if err != nil {
			log.Println("fixture service:", err)
			return err
		}

		fixtures[i] = f
	}

	for _, f := range fixtures {
		err = fs.fr.Update(f)
		if err == storage.ErrFixtureNotFound {
			err = fs.fr.Add(f)
			if err != nil {
				log.Println("fixture service:", err)
				return err
			}
		} else if err != nil {
			log.Println("fixture service:", err)
			return err
		}
	}

	return nil
}

// GetFixturesByGameweek returns all fixtures that take place during gameweek
// with given ID, or returns error otherwise
func (fs *fixtureService) GetFixturesByGameweek(gameweekID int) ([]domain.Fixture, error) {
	fixture := domain.Fixture{GameweekID: gameweekID}

	err := runFixtureValidations(&fixture, gameweekIDBetween1and38)
	if err != nil {
		return []domain.Fixture{}, err
	}

	return fs.fr.GetByGameweek(gameweekID)
}

// GetLiveFixtures returns all fixtures with bonus points not yet confirmed,
// from gameweek with given ID. Returns error on failure.
func (fs *fixtureService) GetLiveFixtures(gameweekID int) ([]domain.Fixture, error) {
	gwFixtures, err := fs.GetFixturesByGameweek(gameweekID)
	if err != nil {
		return []domain.Fixture{}, err
	}

	liveFixtures := make([]domain.Fixture, 0)
	for _, f := range gwFixtures {
		if f.Started && !f.Finished {
			liveFixtures = append(liveFixtures, f)
		}
	}

	return liveFixtures, nil
}

// convertToDomainFixture returns domain.Fixture, consistent with given
// wrapper.Fixture, returns error if it fails to parse fixture's kickoff time
func (fs *fixtureService) convertToDomainFixture(wf wrapper.Fixture) (domain.Fixture, error) {
	clubHome, err := fs.cs.GetClubByID(wf.TeamH)
	if err != nil {
		return domain.Fixture{}, err
	}

	clubAway, err := fs.cs.GetClubByID(wf.TeamA)
	if err != nil {
		return domain.Fixture{}, err
	}

	var kickoffTime time.Time
	if wf.Event != 0 { // if fixture is scheduled
		kickoffTime, err = time.Parse(time.RFC3339, wf.KickoffTime)
		if err != nil {
			log.Println("fixture service:", err)
			return domain.Fixture{}, err
		}
	}

	stats := make(map[string]domain.FixtureStat, len(wf.Stats))
	for _, s := range wf.Stats {
		homePlayersStats := make([]domain.FixtureStatPair, len(s.TeamH))
		awayPlayersStats := make([]domain.FixtureStatPair, len(s.TeamA))

		for _, stat := range s.TeamH {
			homePlayersStats = append(homePlayersStats, domain.FixtureStatPair{
				PlayerID: stat.Element,
				Value:    stat.Value,
			})
		}

		for _, stat := range s.TeamA {
			awayPlayersStats = append(awayPlayersStats, domain.FixtureStatPair{
				PlayerID: stat.Element,
				Value:    stat.Value,
			})
		}

		fixtureStat := domain.FixtureStat{
			Name:             s.Identifier,
			HomePlayersStats: homePlayersStats,
			AwayPlayersStats: awayPlayersStats,
		}

		stats[fixtureStat.Name] = fixtureStat
	}

	fixture := domain.Fixture{
		GameweekID:          wf.Event,
		ID:                  wf.ID,
		ClubHome:            clubHome,
		ClubAway:            clubAway,
		Started:             wf.Started,
		Finished:            wf.Finished,
		FinishedProvisional: wf.FinishedProvisional,
		KickoffTime:         kickoffTime,
		Stats:               stats,
	}

	return fixture, nil
}
