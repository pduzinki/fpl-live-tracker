package club

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/wrapper"
)

// ClubService is an interface for interacting with clubs
type ClubService interface {
	GetClubByID(id int) (domain.Club, error)
}

// clubService implements ClubService interface
type clubService struct {
	clubs domain.ClubRepository
}

// NewClubService creates new instance of ClubService, and fills underlying storage with data from FPL API
func NewClubService(clubRepo domain.ClubRepository, w wrapper.Wrapper) (ClubService, error) {
	cs := &clubService{
		clubs: clubRepo,
	}

	wrapperClubs, err := w.GetClubs()
	if err != nil {
		return nil, err
	}

	for _, wc := range wrapperClubs {
		club := cs.convertToDomainClub(wc)

		err = clubRepo.Add(club)
		if err != nil {
			return nil, err
		}
	}

	return cs, nil
}

// GetClubByID returns domain.Club with given ID, or returns error otherwise
func (cs *clubService) GetClubByID(id int) (domain.Club, error) {
	club := domain.Club{ID: id}

	err := runClubValidations(&club, idBetween1and20)
	if err != nil {
		return domain.Club{}, err
	}

	return cs.clubs.GetByID(id)
}

// convertToDomainClub returns domain.Club object, consistent with given wrapper.Club object
func (cs *clubService) convertToDomainClub(wc wrapper.Club) domain.Club {
	return domain.Club{
		ID:        wc.ID,
		Name:      wc.Name,
		Shortname: wc.Shortname,
	}
}
