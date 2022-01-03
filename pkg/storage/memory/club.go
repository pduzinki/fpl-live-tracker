package memory

import (
	"errors"
	domain "fpl-live-tracker/pkg"
	"sync"
)

var (
	ErrClubNotFound      error = errors.New("storage: club not found")
	ErrClubAlreadyExists error = errors.New("storage: club already exists")
)

//
type clubRepository struct {
	clubs map[int]domain.Club
	sync.Mutex
}

//
func NewClubRepository() domain.ClubRepository {
	return &clubRepository{
		clubs: make(map[int]domain.Club),
	}
}

//
func (cr *clubRepository) Add(club domain.Club) error {
	if _, ok := cr.clubs[club.ID]; ok {
		return ErrClubAlreadyExists
	}

	cr.Lock()
	cr.clubs[club.ID] = club
	cr.Unlock()

	return nil
}

//
func (cr *clubRepository) AddMany(clubs []domain.Club) error {
	for _, club := range clubs {
		err := cr.Add(club)
		if err != nil {
			return err
		}
	}
	return nil
}

//
func (cr *clubRepository) GetByID(id int) (domain.Club, error) {
	if club, ok := cr.clubs[id]; ok {
		return club, nil
	}

	return domain.Club{}, ErrClubNotFound
}