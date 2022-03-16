package wrapper

import (
	"errors"
	"fmt"
)

var ErrReadFailure error = errors.New("failed to read the response")
var ErrUnmarshalFailure error = errors.New("failed to unmarshal data")

type ErrorHttpNotOk struct {
	statusCode int
}

func (err ErrorHttpNotOk) Error() string {
	return fmt.Sprintf("http status not ok: %d\n", err.statusCode)
}

func (err ErrorHttpNotOk) GetHttpStatusCode() int {
	return err.statusCode
}
