package manager

import (
	"fpl-live-tracker/pkg/domain"
	"fpl-live-tracker/pkg/wrapper"
)

type ManagerService interface {
	Update() error
	GetByID(id int) (domain.Manager, error)
}

type managerService struct {
	mr domain.ManagerRepository
	wr wrapper.Wrapper
}

//
func NewManagerService(mr domain.ManagerRepository, wr wrapper.Wrapper) (ManagerService, error) {
	ms := managerService{
		mr: mr,
		wr: wr,
	}

	err := ms.Update()
	if err != nil {
		return nil, err
	}

	return &ms, nil
}

func (ms *managerService) Update() error {
	return nil
}

//
func (ms *managerService) GetByID(id int) (domain.Manager, error) {
	return domain.Manager{}, nil
}

func (ms *managerService) convertToDomainManager() domain.Manager {
	return domain.Manager{}
}
