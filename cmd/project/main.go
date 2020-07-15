package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/ikascrew/client/tool"
)

func main() {

	flag.Parse()
	args := flag.Args()

	l := len(args)
	fmt.Println(args)
	if l < 1 {
		os.Exit(1)
	}

	project := args[0]
	err := tool.CreateProject(project)

	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	os.Exit(0)
}
