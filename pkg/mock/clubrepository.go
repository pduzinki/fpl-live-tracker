package mock

import "fpl-live-tracker/pkg/domain"

type ClubRepository struct {
	AddFn     func(club domain.Club) error
	AddManyFn func(clubs []domain.Club) error
	GetByIDFn func(id int) (domain.Club, error)
}

func (cr *ClubRepository) Add(club domain.Club) error {
	return cr.AddFn(club)
}

func (cr *ClubRepository) AddMany(clubs []domain.Club) error {
	return cr.AddManyFn(clubs)
}

func (cr *ClubRepository) GetByID(id int) (domain.Club, error) {
	return cr.GetByIDFn(id)
}
