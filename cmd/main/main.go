package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/rsmanito/developstoday-test-assessment/app"
	"github.com/rsmanito/developstoday-test-assessment/server"
)

func main() {
	s := server.New()

	app := app.New(s)

	errChan := make(chan error, 1)
	go func() {
		// TODO: move port to config
		if err := app.Run(":3000"); err != nil {
			errChan <- err
		}
	}()

	log.Default().Println("Server is running on port 3000")

	// Capture signals to perform a graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		log.Default().Println("Received a startup error: ", err)
	case sig := <-sigChan:
		log.Default().Println("Received a signal: ", sig)
	}

	if err := app.Shutdown(); err != nil {
		log.Default().Println("Received a shutdown error: ", err)
	} else {
		log.Default().Println("Server shutdown gracefully")
	}
}
