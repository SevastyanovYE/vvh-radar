package main

import (
	"log"

	"github.com/example/vvh-radar/internal/cli"
)

func main() {
	if err := cli.NewRoot().Execute(); err != nil {
		log.Fatal(err)
	}
}
