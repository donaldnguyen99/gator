package main

import (
	"fmt"
	"log"

	"github.com/donaldnguyen99/gator/internal/config"
)

func main() {
	c, err := config.Read()
	if err != nil {
		log.Fatal(err)
	}

	c.SetUser("donald")

	c, err = config.Read()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(*c)
}