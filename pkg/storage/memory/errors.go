package memory

import (
	"errors"
	"fmt"
)

// var ErrRecordAlreadyExists error = errors.New("storage: record already exists") // TODO not sure if that's the best place for those
var ErrRecordNotFound error = errors.New("storage: record not found")

type errorRecordAlreadyExists struct {
	fplID int
}

type ErrorRecordAlreadyExists interface {
	error
	GetFplID() int
}

func (err errorRecordAlreadyExists) Error() string {
	return fmt.Sprintf("storage: record with fplID '%d' already exists\n", err.fplID)
}

func (err errorRecordAlreadyExists) GetFplID() int {
	return err.fplID
}
