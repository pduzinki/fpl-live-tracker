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

	//
	// wrapper.GetClubs()

	return cs, nil
}

func (cs *clubService) GetClubByID(id int) (domain.Club, error) {
	return domain.Club{}, nil
}
