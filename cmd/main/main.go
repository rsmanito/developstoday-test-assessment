package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/rsmanito/developstoday-test-assessment/app"
	"github.com/rsmanito/developstoday-test-assessment/config"
	"github.com/rsmanito/developstoday-test-assessment/server"
	"github.com/rsmanito/developstoday-test-assessment/service"
	"github.com/rsmanito/developstoday-test-assessment/storage"
)

func main() {
	slog.SetDefault(
		slog.New(slog.NewTextHandler(os.Stdout, nil)),
	)

	cfg := config.MustLoad()

	storage := storage.New(cfg)
	service := service.New(storage)
	server := server.New(service)

	app := app.New(server)

	errChan := make(chan error, 1)
	go func() {
		if err := app.Run(":" + cfg.HttpPort); err != nil {
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
