package main

import (
	"log"
	"os"
	"poly-cli/internal/cli"
)

const (
	cmdGenerate = "generate"
)

func main() {
	args := os.Args
	cmd := args[1]

	var err error = nil
	switch cmd {
	case cmdGenerate:
		err = cli.Generate()

	default:
		log.Fatalln("Unknown command: " + cmd)
	}

	if err != nil {
		log.Fatal(err)
	}
}
