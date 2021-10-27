package main

import (
	"fmt"
	"fpl-live-tracker/pkg/storage/memory"
)

func main() {
	fmt.Println("hello there, from fpl-live-tracker!")

	mr, err := memory.NewManagerRepository()
	if err != nil {
		panic(err)
	}
	_ = mr
}
