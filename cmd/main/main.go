package main

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rsmanito/developstoday-test-assessment/app"
	"github.com/rsmanito/developstoday-test-assessment/server"
)

func main() {
	s := server.New()
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)

	app := app.New(s)

	errChan := make(chan error, 1)
	go func() {
		// TODO: move port to config
		if err := app.Run(":3000"); err != nil {
			errChan <- err
		}
	}()

	slog.Info("Server is running", "port", cfg.HttpPort)

	// Capture signals to perform a graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case err := <-errChan:
		slog.Error("Receiver a startup error", "err", err)
	case sig := <-sigChan:
		slog.Error("Received signal", "sig", sig)
	}

	if err := app.Shutdown(); err != nil {
		slog.Error("Received shutdown error", "err", err)
	} else {
		slog.Info("Service shutdown gracefully")
	}
}
