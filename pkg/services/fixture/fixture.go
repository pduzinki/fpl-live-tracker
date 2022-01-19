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
func NewFixtureService(fr domain.FixtureRepository, cs club.ClubService, w wrapper.Wrapper) FixtureService {
	return &fixtureService{
		fr:      fr,
		cs:      cs,
		wrapper: w,
	}
}

//
func (fs *fixtureService) Update() error {
	wrapperFixtures, err := fs.wrapper.GetFixtures()
	if err != nil {
		return err
	}

	fixtures := make([]domain.Fixture, len(wrapperFixtures))
	for i, wf := range wrapperFixtures {
		clubHome, _ := fs.cs.GetClubByID(wf.TeamH)
		clubAway, _ := fs.cs.GetClubByID(wf.TeamA)

		var kickoffTime time.Time
		if wf.Event != 0 { // if fixture is scheduled
			kickoffTime, err = time.Parse(time.RFC3339, wf.KickoffTime)
			if err != nil {
				log.Println("fixture service:", err)
				return err
			}
		}

		stats := make([]domain.FixtureStat, 0)
		for _, s := range wf.Stats {
			tmpH := make([]domain.FixtureStatValue, 0)
			tmpA := make([]domain.FixtureStatValue, 0)

			for _, item := range s.TeamH {
				tmpH = append(tmpH, domain.FixtureStatValue{
					PlayerID: item.Element,
					Value:    item.Value,
				})
			}

			for _, item := range s.TeamA {
				tmpA = append(tmpA, domain.FixtureStatValue{
					PlayerID: item.Element,
					Value:    item.Value,
				})
			}

			tmp := domain.FixtureStat{
				Name:             s.Identifier,
				HomePlayersStats: tmpH,
				AwayPlayersStats: tmpA,
			}

			stats = append(stats, tmp)
		}

		fixtures[i] = domain.Fixture{
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
