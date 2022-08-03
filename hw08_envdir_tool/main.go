package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args

	environment, err := ReadDir(args[0])
	if err != nil {
		log.Fatal(err.Error())
	}

	fmt.Println(RunCmd(args[1:], environment))
}
