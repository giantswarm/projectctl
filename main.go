package main

import (
	"log"

	"giantswarm.io/projectctl/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatal(err)
	}

}
