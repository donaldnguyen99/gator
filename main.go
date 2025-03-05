package main

import (
	"log"

	"github.com/donaldnguyen99/gator/internal/cli"
)

func main() {
	gator_cli := cli.NewCLI("gator")
	if err := gator_cli.Run(); err != nil {
		log.Fatal(err)
	}
}