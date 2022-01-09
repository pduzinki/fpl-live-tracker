package wrapper

import (
	"errors"
	"fmt"
)

var ErrReadFailure error = errors.New("failed to read the response")
var ErrUnmarshalFailure error = errors.New("failed to unmarshal data")

type errorHttpNotOk struct {
	statusCode int
}

type ErrorHttpNotOk interface {
	error
	GetHttpStatusCode() int
}

func (err errorHttpNotOk) Error() string {
	return fmt.Sprintf("http status not ok: %d\n", err.statusCode)
}

func (err errorHttpNotOk) GetHttpStatusCode() int {
	return err.statusCode
}
