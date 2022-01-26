package fixture

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"log"
	"time"
)

type FixtureService interface {
	Update() error
	GetFixturesByGameweek(gameweekID int) ([]domain.Fixture, error)
	GetLiveFixtures(gameweekID int) ([]domain.Fixture, error)
}

type fixtureService struct {
	fr      domain.FixtureRepository
	cs      club.ClubService
	wrapper wrapper.Wrapper
}

//
func NewFixtureService(fr domain.FixtureRepository, cs club.ClubService, w wrapper.Wrapper) (FixtureService, error) {
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

//
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

//
func (fs *fixtureService) GetFixturesByGameweek(gameweekID int) ([]domain.Fixture, error) {
	fixture := domain.Fixture{GameweekID: gameweekID}

	err := runFixtureValidations(&fixture, gameweekIDBetween1and38)
	if err != nil {
		return []domain.Fixture{}, err
	}

	return fs.fr.GetByGameweek(gameweekID)
}

//
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
		homePlayersStats := make([]domain.FixtureStatValue, len(s.TeamH))
		awayPlayersStats := make([]domain.FixtureStatValue, len(s.TeamA))

		for _, stat := range s.TeamH {
			homePlayersStats = append(homePlayersStats, domain.FixtureStatValue{
				PlayerID: stat.Element,
				Value:    stat.Value,
			})
		}

		for _, stat := range s.TeamA {
			awayPlayersStats = append(awayPlayersStats, domain.FixtureStatValue{
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
