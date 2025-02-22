package server

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/rsmanito/developstoday-test-assessment/models"
)

type Service interface {
	GetAllCats(context.Context) ([]models.Cat, error)
	CreateCat(context.Context, models.CreateCatRequest) (models.Cat, error)
	GetCat(context.Context, int32) (models.Cat, error)
	UpdateCatSalary(context.Context, models.UpdateCatSalaryRequest, int32) (models.Cat, error)
	DeleteCat(context.Context, int32) error
}

type Server struct {
	service Service
	R       *fiber.App
}

// New returns a new Server.
func New(service Service) Server {
	server := Server{
		service: service,
		R: fiber.New(
			fiber.Config{
				StructValidator: &models.StructValidator{Validator: validator.New()},
			},
		),
	}

	server.R.Use(LoggerMiddleware())

	server.registerRoutes()

	return server
}

// registerRoutes registers routes for the Server.
func (s *Server) registerRoutes() {
	// Cats group
	cats := s.R.Group("/cats")
	{
		cats.Get("/", s.handleGetCats)
		cats.Post("/", s.handleCreateCat)
		cats.Get("/:id", s.handleGetSingleCat)
		cats.Patch("/:id", s.handleUpdateCatSalary)
		cats.Delete("/:id", s.handleDeleteCat)
	}
}

func (s *Server) handleGetCats(c fiber.Ctx) error {
	res, err := s.service.GetAllCats(c.Context())
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"cats": res})
}

func (s *Server) handleCreateCat(c fiber.Ctx) error {
	r := models.CreateCatRequest{}

	if err := c.Bind().JSON(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := s.service.CreateCat(c.Context(), r)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (s *Server) handleGetSingleCat(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	res, err := s.service.GetCat(c.Context(), int32(id))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) handleUpdateCatSalary(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var r models.UpdateCatSalaryRequest

	if err := c.Bind().JSON(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := s.service.UpdateCatSalary(c.Context(), r, int32(id))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) handleDeleteCat(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = s.service.DeleteCat(c.Context(), int32(id))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func handleError(c fiber.Ctx, err error) error {
	var customErr *models.Err
	if errors.As(err, &customErr) {
		return c.Status(customErr.Code).JSON(fiber.Map{
			"error": customErr.Msg,
		})
	}

	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error": err.Error(),
	})
}
