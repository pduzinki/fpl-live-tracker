package fixture

import (
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/services/club"
	"fpl-live-tracker/pkg/wrapper"
)

type FixtureService interface {
	Update(gameweekID int) error
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
func (fs *fixtureService) Update(gameweekID int) error {
	wrapperFixtures, err := fs.wrapper.GetFixtures(gameweekID)
	if err != nil {
		// return []domain.Fixture{}, err
		return err
	}
	_ = wrapperFixtures

	fixtures := make([]domain.Fixture, 0)
	for _, wf := range wrapperFixtures {
		clubAway, _ := fs.cs.GetClubByID(wf.TeamA)
		clubHome, _ := fs.cs.GetClubByID(wf.TeamH)

		fixtures = append(fixtures, domain.Fixture{
			GameweekID: wf.Event,
			ID:         wf.ID,
			ClubAway:   clubAway,
			ClubHome:   clubHome,
		})
	}

	err = fs.fr.AddMany(fixtures)
	if err != nil {
		panic(err)
	}

	// return fixtures, nil
	return nil
}

//
func (fs *fixtureService) GetFixturesByGameweek(gameweekID int) ([]domain.Fixture, error) {
	// TODO add validations
	return fs.fr.GetByGameweek(gameweekID)
}
