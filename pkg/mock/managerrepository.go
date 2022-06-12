package mock

import "fpl-live-tracker/pkg/domain"

type ManagerRepository struct {
	AddFn      func(manager domain.Manager) error
	AddManyFn  func(managers []domain.Manager) error
	UpdateFn   func(manager domain.Manager) error
	GetByIDFn  func(id int) (domain.Manager, error)
	GetCountFn func() (int, error)
}

func (mr *ManagerRepository) Add(manager domain.Manager) error {
	return mr.AddFn(manager)
}

func (mr *ManagerRepository) AddMany(managers []domain.Manager) error {
	return mr.AddManyFn(managers)
}

func (mr *ManagerRepository) Update(manager domain.Manager) error {
	return mr.UpdateFn(manager)
}

func (mr *ManagerRepository) GetByID(id int) (domain.Manager, error) {
	return mr.GetByIDFn(id)
}

func (mr *ManagerRepository) GetCount() (int, error) {
	return mr.GetCountFn()
}
