package main

import (
	"fmt"

	"github.com/google/uuid"
)

func main() {
	n := 10
	for i := 0; i < n; i++ {
		fmt.Println(uuid.New().String())
	}
}
