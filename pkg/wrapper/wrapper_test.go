package wrapper

import (
	"fmt"
	"testing"
)

func TestDummy(t *testing.T) {
	w := NewWrapper(DefaultURL)

	manager, err := w.GetManager(1239)
	if err != nil {
		panic(err)
	}

	fmt.Println(manager)
}
