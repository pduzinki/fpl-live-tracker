package storage

import "errors"

// club errors
var (
	ErrClubNotFound      error = errors.New("storage: club not found")
	ErrClubAlreadyExists error = errors.New("storage: club already exists")
)

// fixture errors
var (
	ErrFixtureNotFound      error = errors.New("storage: fixture not found")
	ErrFixtureAlreadyExists error = errors.New("storage: fixture already exists")
)

// player errors
var (
	ErrPlayerNotFound      error = errors.New("storage: player not found")
	ErrPlayerAlreadyExists error = errors.New("storage: player already exists")
)

// manager errors
var (
	ErrManagerAlreadyExists error = errors.New("storage: manager already exists")
	ErrManagerNotFound      error = errors.New("storage: manager not found")
)

// team errors
var (
	ErrTeamAlreadyExists error = errors.New("storage: team already exists")
	ErrTeamNotFound      error = errors.New("storage: team not found")
)
