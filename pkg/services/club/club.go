package club

import (
	domain "fpl-live-tracker/pkg"
	"fpl-live-tracker/pkg/wrapper"
)

type ClubService interface {
	// Add(domain.Club) is not needed, all clubs are added when service is created
	GetClubByID(id int) (domain.Club, error)
}

type clubService struct {
	clubs domain.ClubRepository
}

func NewClubService(clubRepo domain.ClubRepository, wrapper wrapper.Wrapper) (ClubService, error) {
	cs := &clubService{
		clubs: clubRepo,
	}

	wrapperClubs, err := wrapper.GetClubs() // TODO to add http retries would be nice
	if err != nil {
		return nil, err
	}

	for _, wc := range wrapperClubs {
		club := domain.Club{
			ID:        wc.ID,
			Name:      wc.Name,
			Shortname: wc.Shortname,
		}

		err = clubRepo.Add(club)
		if err != nil {
			return nil, err
		}
	}

	return cs, nil
}

func (cs *clubService) GetClubByID(id int) (domain.Club, error) {
	// TODO add validations

	return cs.clubs.GetByID(id)
}
