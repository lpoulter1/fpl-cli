package main

import (
	"log"

	"github.com/lpt10/fpl-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
