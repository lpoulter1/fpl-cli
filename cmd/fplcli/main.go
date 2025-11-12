package main

import (
	"log"

	"github.com/lpoulter1/fpl-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
