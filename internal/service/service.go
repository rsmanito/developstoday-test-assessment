package service

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rsmanito/developstoday-test-assessment/internal/storage/postgres"
)

// CatStorage controls the cat storage.
type CatStorage interface {
	GetAllCats(ctx context.Context) ([]postgres.Cat, error)
	CreateCat(ctx context.Context, params postgres.CreateCatParams) (postgres.Cat, error)
	GetCat(ctx context.Context, id int32) (postgres.Cat, error)
	UpdateCatSalary(ctx context.Context, params postgres.UpdateCatSalaryParams) (postgres.Cat, error)
	DeleteCat(ctx context.Context, id int32) (int64, error)
}

// MissionStorage controls the mission storage.
type MissionStorage interface {
	CreateMission(ctx context.Context) (postgres.Mission, error)
	GetAllMissions(ctx context.Context) ([]postgres.Mission, error)
	GetMission(ctx context.Context, id int32) (postgres.Mission, error)
	DeleteMission(ctx context.Context, id int32) (int64, error)
	AssignCat(ctx context.Context, params postgres.AssignCatParams) (postgres.Mission, error)
	GetCatMission(ctx context.Context, catID pgtype.Int4) (postgres.Mission, error)
	CompleteMission(ctx context.Context, id int32) (postgres.Mission, error)
	GetMissionByTargetID(ctx context.Context, id int32) (postgres.GetMissionByTargetIDRow, error)
}

// TargetStorage controls the target storage.
type TargetStorage interface {
	GetTarget(ctx context.Context, id int32) (postgres.Target, error)
	GetMissionTargets(ctx context.Context, missionID int32) ([]postgres.Target, error)
	CreateTarget(ctx context.Context, params postgres.CreateTargetParams) (postgres.Target, error)
	DeleteTarget(ctx context.Context, id int32) (int64, error)
	UpdateTargetNotes(ctx context.Context, params postgres.UpdateTargetNotesParams) (postgres.Target, error)
	CompleteTarget(ctx context.Context, id int32) (postgres.Target, error)
}

// TransactionalStorage controls the transactional storage.
type TransactionalStorage interface {
	WithTx(tx pgx.Tx) *postgres.Queries
	Begin(ctx context.Context) (pgx.Tx, error)
}

type Service struct {
	catStorage     CatStorage
	missionStorage MissionStorage
	targetStorage  TargetStorage
	txStorage      TransactionalStorage
}

// New returns a new Service.
func NewService(cs CatStorage, ms MissionStorage, ts TargetStorage, txs TransactionalStorage) Service {
	return Service{
		catStorage:     cs,
		missionStorage: ms,
		targetStorage:  ts,
		txStorage:      txs,
	}
}
