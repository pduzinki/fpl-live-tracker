package mock

import "fpl-live-tracker/pkg/domain"

type ManagerRepository struct {
	AddFn        func(manager domain.Manager) error
	AddManyFn    func(managers []domain.Manager) error
	UpdateInfoFn func(managerID int, info domain.ManagerInfo) error
	UpdateTeamFn func(managerID int, team domain.Team) error
	GetByIDFn    func(id int) (domain.Manager, error)
}

func (mr *ManagerRepository) Add(manager domain.Manager) error {
	return mr.AddFn(manager)
}

func (mr *ManagerRepository) AddMany(managers []domain.Manager) error {
	return mr.AddManyFn(managers)
}

func (mr *ManagerRepository) UpdateInfo(managerID int, info domain.ManagerInfo) error {
	return mr.UpdateInfoFn(managerID, info)
}

func (mr *ManagerRepository) UpdateTeam(managerID int, team domain.Team) error {
	return mr.UpdateTeamFn(managerID, team)
}

func (mr *ManagerRepository) GetByID(id int) (domain.Manager, error) {
	return mr.GetByIDFn(id)
}
