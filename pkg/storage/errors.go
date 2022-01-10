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

// manager errors
var (
	ErrManagerAlreadyExists error = errors.New("storage: manager already exists")
	ErrManagerNotFound      error = errors.New("storage: manager not found")
)
