package memory

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/storage"
	"sync"
)

// clubRepository implements domain.ClubRepository interface
type clubRepository struct {
	clubs map[int]domain.Club
	sync.Mutex
}

// NewClubRepository returns new instance of domain.ClubRepository
func NewClubRepository() domain.ClubRepository {
	return &clubRepository{
		clubs: make(map[int]domain.Club),
	}
}

// Add saves given club into memory storage, returns error on failure
func (cr *clubRepository) Add(club domain.Club) error {
	if _, ok := cr.clubs[club.ID]; ok {
		return storage.ErrClubAlreadyExists
	}

	cr.Lock()
	cr.clubs[club.ID] = club
	cr.Unlock()

	return nil
}

// AddMany saves given clubs into memory storage, returns error on failure
func (cr *clubRepository) AddMany(clubs []domain.Club) error {
	for _, club := range clubs {
		err := cr.Add(club)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetByID returns club with given id, or error otherwise
func (cr *clubRepository) GetByID(id int) (domain.Club, error) {
	if club, ok := cr.clubs[id]; ok {
		return club, nil
	}

	return domain.Club{}, storage.ErrClubNotFound
}
