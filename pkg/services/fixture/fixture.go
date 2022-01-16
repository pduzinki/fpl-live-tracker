package fixture

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/storage"
	"fpl-live-tracker/pkg/wrapper"
	"log"
)

type FixtureService interface {
	Update() error
	GetFixturesByGameweek(gameweekID int) ([]domain.Fixture, error)
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

		fixtures[i] = domain.Fixture{
			GameweekID: wf.Event,
			ID:         wf.ID,
			ClubHome:   clubHome,
			ClubAway:   clubAway,
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
