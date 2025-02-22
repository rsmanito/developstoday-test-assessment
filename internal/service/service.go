package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rsmanito/developstoday-test-assessment/internal/storage/postgres"
)

type Storage interface {
	GetAllCats(context.Context) ([]postgres.Cat, error)
	CreateCat(context.Context, postgres.CreateCatParams) (postgres.Cat, error)
	GetCat(context.Context, int32) (postgres.Cat, error)
	UpdateCatSalary(context.Context, postgres.UpdateCatSalaryParams) (postgres.Cat, error)
	DeleteCat(context.Context, int32) (int64, error)

	CreateMission(context.Context) (postgres.Mission, error)
	GetAllMissions(context.Context) ([]postgres.Mission, error)
	GetMission(context.Context, int32) (postgres.Mission, error)
	DeleteMission(context.Context, int32) (int64, error)
	AssignCat(context.Context, postgres.AssignCatParams) (postgres.Mission, error)
	GetCatMission(context.Context, pgtype.Int4) (postgres.Mission, error)
	CompleteMission(context.Context, int32) (postgres.Mission, error)
	GetMissionByTargetID(context.Context, int32) (postgres.GetMissionByTargetIDRow, error)

	GetTarget(context.Context, int32) (postgres.Target, error)
	GetMissionTargets(context.Context, int32) ([]postgres.Target, error)
	CreateTarget(context.Context, postgres.CreateTargetParams) (postgres.Target, error)
	DeleteTarget(context.Context, int32) (int64, error)
	UpdateTargetNotes(context.Context, postgres.UpdateTargetNotesParams) (postgres.Target, error)
	CompleteTarget(context.Context, int32) (postgres.Target, error)

	WithTx(pgx.Tx) *postgres.Queries
	Begin(context.Context) (pgx.Tx, error)
}

type Service struct {
	st Storage
}

// New returns a new Service.
func New(storage Storage) Service {
	return Service{
		st: storage,
	}
}
