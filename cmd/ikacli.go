package main

import (
	"log"
	"os"

	"github.com/ikascrew/client"
)

const VERSION = "0.0.0"

func main() {

	err := client.Start()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
