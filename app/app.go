package app

import (
	"context"
	"time"

	"github.com/rsmanito/developstoday-test-assessment/server"
)

type App struct {
	httpServer *server.Server
}

// New returns a new App.
func New(s *server.Server) *App {
	return &App{
		httpServer: s,
	}
}

// Run starts the application.
//
// Returns an error if something goes wrong.
func (a *App) Run(addr string) error {
	err := a.httpServer.R.Listen(addr)
	if err != nil {
		return err
	}
	return nil
}

// Shutdown performs a graceful shutdown after a timeout.
//
// Returns an error if something goes wrong.
func (a *App) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.httpServer.R.ShutdownWithContext(ctx); err != nil {
		return err
	}
	return nil
}
