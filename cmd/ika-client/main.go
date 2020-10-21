package main

import (
	"log"
	"os"

	"github.com/ikascrew/client"
	"github.com/ikascrew/client/config"
)

func main() {

	err := client.Start(
		config.Controller(config.ControllerTypeJoyCon),
		//config.UsePowermate(),
		//config.Windows(),
	)
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
