package memory

import (
	"errors"
	"fmt"
)

// TODO move to storage/errors.go
// var ErrRecordAlreadyExists error = errors.New("storage: record already exists")
var ErrManagerNotFound error = errors.New("storage: manager not found")

type errManagerAlreadyExists struct {
	fplID int
}

type ErrManagerAlreadyExists interface {
	error
	GetFplID() int
}

func (err errManagerAlreadyExists) Error() string {
	return fmt.Sprintf("storage: manager with fplID '%d' already exists\n", err.fplID)
}

func (err errManagerAlreadyExists) GetFplID() int {
	return err.fplID
}
