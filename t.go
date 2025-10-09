package main

import (
	"fmt"

	"github.com/raulaguila/go-api/pkg/packhub"
)

func main() {
	fmt.Println(packhub.RandomFloat64(0.1, 0.7, 2))
	fmt.Println(packhub.RandomFloat64(5, 7, 2))
}
