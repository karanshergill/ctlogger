package main

import (
	"log"

	"github.com/karanshergill/ctlogger/pkg/daemon"
	"github.com/karanshergill/ctlogger/pkg/runner"
)

func main() {
	options, err := runner.ParseOptions()
	if err != nil {
		log.Fatalf("Error parsing options: %v", err)
	}

	// Handle daemon mode
	if options.Daemon {
		err := daemon.Daemonize(options.PidFile)
		if err != nil {
			log.Fatalf("Error starting daemon: %v", err)
		}
	}

	runner, err := runner.NewRunner(options)
	if err != nil {
		log.Fatalf("Error creating runner: %v", err)
	}

	runner.Run()
}
