package main

import (
	"log"
	"os"
	"poly-cli/internal/cli"
)

const (
	cmdGenerate  = "generate"
	cmdOpen      = "open"
	cmdRun       = "run"
	cmdDevServer = "dev-server"
)

func main() {
	args := os.Args
	cmd := args[1]

	var err error = nil
	switch cmd {
	case cmdGenerate:
		err = cli.Generate()
	case cmdOpen:
		err = cli.Open(args[2:]...)
	case cmdDevServer:
		err = cli.StartDevServer()

	default:
		log.Fatalln("Unknown command: " + cmd)
	}

	if err != nil {
		log.Fatal(err)
	}
}
