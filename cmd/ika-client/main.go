package main

import (
	"log"
	"os"

	"github.com/ikascrew/client"
	"github.com/ikascrew/client/config"
)

const VERSION = "0.0.0"

func main() {

	err := client.Start(
		config.Controller(config.ControllerTypeJoyCon),
		//config.UsePowermate(),
	)
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	os.Exit(0)
}
