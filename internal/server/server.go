package server

import (
	"context"
	"errors"
	"strconv"

	"github.com/go-playground/validator"
	"github.com/gofiber/fiber/v3"
	"github.com/rsmanito/developstoday-test-assessment/internal/models"
)

// CatService controls the cat service.
type CatService interface {
	GetAllCats(ctx context.Context) ([]models.Cat, error)
	CreateCat(ctx context.Context, req models.CreateCatRequest) (models.Cat, error)
	GetCat(ctx context.Context, id int32) (models.Cat, error)
	UpdateCatSalary(ctx context.Context, req models.UpdateCatSalaryRequest, id int32) (models.Cat, error)
	DeleteCat(ctx context.Context, id int32) error
}

// MissionService controls the mission service.
type MissionService interface {
	GetAllMissions(ctx context.Context) ([]models.Mission, error)
	CreateMission(ctx context.Context, req models.CreateMissionRequest) (models.Mission, error)
	GetMission(ctx context.Context, id int32) (models.Mission, error)
	AssignCatToMission(ctx context.Context, missionID int32, assignee int32) (models.Mission, error)
	CompleteMission(ctx context.Context, id int32) (models.Mission, error)
	DeleteMission(ctx context.Context, id int32) error
	AddTarget(ctx context.Context, missionID int32, req models.CreateTargetRequest) (models.Mission, error)
}

// TargetService controls the target service.
type TargetService interface {
	DeleteTarget(ctx context.Context, id int32) error
	UpdateTargetNotes(ctx context.Context, id int32, notes string) (models.Target, error)
	CompleteTarget(ctx context.Context, id int32) (models.Target, error)
}

type Server struct {
	catService     CatService
	missionService MissionService
	targetService  TargetService
	R              *fiber.App
}

// New returns a new Server.
func New(cs CatService, ms MissionService, ts TargetService) Server {
	server := Server{
		catService:     cs,
		missionService: ms,
		targetService:  ts,
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

	// Missions group
	missions := s.R.Group("/missions")
	{
		missions.Post("/", s.handleCreateMission)
		missions.Get("/", s.handleGetMissions)

		withId := missions.Group("/:id")
		{
			withId.Get("/", s.handleGetSingleMission)
			withId.Delete("/", s.handleDeleteMission)
			withId.Patch("/assign", s.handleAssignCat)
			withId.Patch("/complete", s.handleCompleteMission)
		}

		targets := withId.Group("/targets")
		{
			targets.Post("/", s.handleAddTarget)
			targets.Delete("/:targetId", s.handleDeleteTarget)
			targets.Patch("/:targetId/notes", s.handleUpdateTargetNotes)
			targets.Patch("/:targetId/complete", s.handleCompleteTarget)
		}
	}
}

func (s *Server) handleGetCats(c fiber.Ctx) error {
	res, err := s.catService.GetAllCats(c.Context())
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

	res, err := s.catService.CreateCat(c.Context(), r)
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

	res, err := s.catService.GetCat(c.Context(), int32(id))
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

	res, err := s.catService.UpdateCatSalary(c.Context(), r, int32(id))
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

	err = s.catService.DeleteCat(c.Context(), int32(id))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (s *Server) handleCreateMission(c fiber.Ctx) error {
	var r models.CreateMissionRequest

	if err := c.Bind().JSON(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := s.missionService.CreateMission(c.Context(), r)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (s *Server) handleGetMissions(c fiber.Ctx) error {
	res, err := s.missionService.GetAllMissions(c.Context())
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) handleGetSingleMission(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	res, err := s.missionService.GetMission(c.Context(), int32(id))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) handleDeleteMission(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = s.missionService.DeleteMission(c.Context(), int32(id))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{})
}

func (s *Server) handleAssignCat(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var r models.AssignCatRequest

	if err := c.Bind().JSON(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := s.missionService.AssignCatToMission(c.Context(), int32(id), r.Assignee)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) handleCompleteMission(c fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	res, err := s.missionService.CompleteMission(c.Context(), int32(id))
	if err != nil {
		return handleError(c, err)
	}
	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) handleAddTarget(c fiber.Ctx) error {
	missionId, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var r models.CreateTargetRequest

	if err := c.Bind().JSON(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := s.missionService.AddTarget(c.Context(), int32(missionId), r)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusCreated).JSON(res)
}

func (s *Server) handleDeleteTarget(c fiber.Ctx) error {
	targetId, err := strconv.Atoi(c.Params("targetId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	err = s.targetService.DeleteTarget(c.Context(), int32(targetId))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusNoContent).JSON(fiber.Map{})
}

func (s *Server) handleUpdateTargetNotes(c fiber.Ctx) error {
	targetId, err := strconv.Atoi(c.Params("targetId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	var r models.UpdateTargetNotesRequest

	if err := c.Bind().JSON(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	res, err := s.targetService.UpdateTargetNotes(c.Context(), int32(targetId), r.Notes)
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
}

func (s *Server) handleCompleteTarget(c fiber.Ctx) error {
	targetId, err := strconv.Atoi(c.Params("targetId"))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}

	res, err := s.targetService.CompleteTarget(c.Context(), int32(targetId))
	if err != nil {
		return handleError(c, err)
	}

	return c.Status(fiber.StatusOK).JSON(res)
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
