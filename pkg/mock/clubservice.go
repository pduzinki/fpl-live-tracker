package mock

import "fpl-live-tracker/pkg/domain"

type ClubService struct {
	GetClubByIDFn func(int) (domain.Club, error)
}

func (cs *ClubService) GetClubByID(id int) (domain.Club, error) {
	return cs.GetClubByIDFn(id)
}
