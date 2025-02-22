package server

import (
	"github.com/gofiber/fiber/v3"
)

type (
	Service interface{}

	Server struct {
		service Service
		R       *fiber.App
	}
)

// New returns a new Server.
func New() *Server {
	server := &Server{
		R: fiber.New(),
	}

	server.registerRoutes()

	return server
}

// registerRoutes registers routes for the Server.
func (s *Server) registerRoutes() {
	// Cats group
	cats := s.R.Group("/cats")
	{
		cats.Get("/", handleGetCats)
		cats.Post("/", handleCreateCat)
		cats.Get("/:id", handleGetSingleCat)
		cats.Patch("/:id", handleUpdateCatSalary)
		cats.Delete("/:id", handleDeleteCat)
	}
}

func handleGetCats(c fiber.Ctx) error {
	return nil
}

func handleCreateCat(c fiber.Ctx) error {
	return nil
}

func handleGetSingleCat(c fiber.Ctx) error {
	return nil
}

func handleUpdateCatSalary(c fiber.Ctx) error {
	return nil
}

func handleDeleteCat(c fiber.Ctx) error {
	return nil
}
