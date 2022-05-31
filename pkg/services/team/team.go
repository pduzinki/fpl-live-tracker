package team

import (
	"fpl-live-tracker/pkg/domain"
)

// TeamService is an interface for interacting with teams
type TeamService interface {
	UpdateTeams() error
	UpdatePoints() error
	GetByID(id int) (domain.Team, error)
}

//
type teamService struct {
	// wr wrapper.Wrapper
}

//
func NewTeamService() TeamService {
	// TODO implement this
	return &teamService{}
}

//
func (ts *teamService) UpdateTeams() error {
	// TODO implement this
	return nil
}

//
func (ts *teamService) UpdatePoints() error {
	// TODO implement this
	return nil
}

//
func (ts *teamService) GetByID(id int) (domain.Team, error) {
	// TODO implement this
	return domain.Team{}, nil
}
