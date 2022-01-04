package storage

import "errors"

var (
	ErrClubNotFound      error = errors.New("storage: club not found")
	ErrClubAlreadyExists error = errors.New("storage: club already exists")
)

var (
	ErrFixtureNotFound      error = errors.New("storage: fixture not found")
	ErrFixtureAlreadyExists error = errors.New("storage: fixture already exists")
)
