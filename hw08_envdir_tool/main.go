package main

import (
	"log"
	"os"
)

func main() {
	args := os.Args

	environment, err := ReadDir(args[1])
	if err != nil {
		log.Fatal(err.Error())
	}

	os.Exit(RunCmd(args[2:], environment))
}
