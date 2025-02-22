package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/rsmanito/developstoday-test-assessment/models"
	"github.com/rsmanito/developstoday-test-assessment/storage/postgres"
)

func (s Service) GetAllMissions(ctx context.Context) ([]models.Mission, error) {
	log := slog.With(
		slog.String("op", "service.GetAllMissions"),
	)

	log.Debug("Fetching all missions")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := s.st.GetAllMissions(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return make([]models.Mission, 0), models.ErrTimeoutExceeded
		}
		slog.Error("Failed to get missions", "err", err)
		return make([]models.Mission, 0), err
	}

	missions := make([]models.Mission, len(res))
	for i, m := range res {
		missions[i] = sqlcMissionToModel(m)
	}

	log.Debug("Fetched all missions", "res", missions)

	return missions, nil
}

func (s Service) CreateMission(ctx context.Context, req models.CreateMissionRequest) (models.Mission, error) {
	if len(req.Targets) == 0 || len(req.Targets) > 3 {
		return models.Mission{}, models.NewError(http.StatusUnprocessableEntity, "incorrect number of targets (1..3)")
	}
	for _, target := range req.Targets {
		if utf8.RuneCountInString(target.Notes) > 256 {
			return models.Mission{}, models.NewError(http.StatusUnprocessableEntity, "incorrect number of targets (0..256)")
		}
	}

	log := slog.With(
		slog.String("op", "service.CreateMission"),
		slog.Any("req", req),
	)

	log.Debug("Creating a mission")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	tx, err := s.st.Begin(ctx)
	if err != nil {
		slog.Error("Failed to begin transaction", "err", err)
		return models.Mission{}, err
	}
	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	// Create blank mission.
	withTx := s.st.WithTx(tx)

	mission, err := withTx.CreateMission(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		slog.Error("Failed to create mission", "err", err)
		return models.Mission{}, errors.New("failed to create mission")
	}

	log.Debug("Created mission", "id", mission.ID)

	// Create targets for mission.
	for _, t := range req.Targets {
		_, err = withTx.CreateTarget(ctx, postgres.CreateTargetParams{
			Mission: mission.ID,
			Name:    t.Name,
			Country: t.Country,
			Notes:   t.Notes,
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return models.Mission{}, models.ErrTimeoutExceeded
			}
			log.Error("Failed to create target", "err", err)
			return models.Mission{}, errors.New("failed to create mission")
		}
	}

	return sqlcMissionToModel(mission), nil
}

func (s Service) GetMission(ctx context.Context, id int32) (models.Mission, error) {
	log := slog.With(
		slog.String("op", "service.GetMission"),
		slog.Any("id", id),
	)

	log.Debug("Fetching mission")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get mission.
	res, err := s.st.GetMission(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Mission{}, models.ErrNotFound
		}
		slog.Error("Failed to get mission", "err", err)
		return models.Mission{}, errors.New("failed to get mission")
	}

	// Get mission targets.
	targets, err := s.st.GetMissionTargets(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Mission{}, models.ErrNotFound
		}
		slog.Error("Failed to get targets", "err", err)
		return models.Mission{}, errors.New("failed to get mission")
	}

	// Map to model.
	mission := sqlcMissionToModel(res)
	for _, t := range targets {
		mission.Targets = append(mission.Targets, sqlcTargetToModel(t))
	}
	return mission, nil
}

func (s Service) AssignCatToMission(ctx context.Context, mission, assignee int32) (models.Mission, error) {
	log := slog.With(
		slog.String("op", "service.AssignCatToMission"),
		slog.Any("missionId", mission),
		slog.Any("assignee", assignee),
	)

	log.Debug("Assigning cat to mission")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Check if cat exists.
	_, err := s.GetCat(ctx, assignee)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Mission{}, models.NewError(http.StatusNotFound, "cat not found")
		}
		slog.Error("Failed to get cat", "err", err)
		return models.Mission{}, errors.New("cat not found")
	}
	log.Debug("Cat exists")

	// Check if mission exists.
	_, err = s.GetMission(ctx, mission)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Mission{}, models.NewError(http.StatusNotFound, "mission not found")
		}
		slog.Error("failed to get mission", "err", err)
		return models.Mission{}, errors.New("mission not found")
	}

	// Get Current cat mission.
	var lastCatMission postgres.Mission
	lastCatMission, err = s.st.GetCatMission(ctx, pgtype.Int4{Int32: assignee, Valid: true})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			if errors.Is(err, context.DeadlineExceeded) {
				return models.Mission{}, models.ErrTimeoutExceeded
			}
			slog.Error("Failed to get current cat mission", "err", err)
			return models.Mission{}, errors.New("failed to assign cat")
	}

	// Check if cat has active mission.
	if lastCatMission.ID == mission {
		log.Debug("Tried to assign to current mission. Returning")
		return sqlcMissionToModel(lastCatMission), nil
	}

	tx, err := s.st.Begin(ctx)
	if err != nil {
		slog.Error("Failed to begin transaction", "err", err)
		return models.Mission{}, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback(ctx)
		} else {
			_ = tx.Commit(ctx)
		}
	}()

	withTx := s.st.WithTx(tx)

	// Assign to new mission.
	newMission, err := withTx.AssignCat(ctx, postgres.AssignCatParams{
		ID:       mission,
		Assignee: pgtype.Int4{Int32: assignee, Valid: true},
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		slog.Error("failed to assign to new mission", "err", err)
		return models.Mission{}, errors.New("failed to assign cat")
	}

	if lastCatMission.ID != 0 {
		log.Debug("Cat has active mission")
		// Unassign from last mission.
		_, err = withTx.AssignCat(ctx, postgres.AssignCatParams{
			ID:       lastCatMission.ID,
			Assignee: pgtype.Int4{Valid: false},
		})
		if err != nil {
			if errors.Is(err, context.DeadlineExceeded) {
				return models.Mission{}, models.ErrTimeoutExceeded
			}
			slog.Error("failed to unassign from last mission", "err", err, "lastCatMission", lastCatMission.ID)
			return models.Mission{}, errors.New("failed to assign cat")
		}

		log.Debug("Unassigned from last mission", "lastCatMission", lastCatMission.ID)
	}

	log.Debug("Assigned to new mission")


	return sqlcMissionToModel(newMission), nil
}

func (s Service) CompleteMission(ctx context.Context, mission int32) (models.Mission, error) {
	log := slog.With(
		slog.String("op", "service.CompleteMission"),
		slog.Any("id", mission),
	)

	log.Debug("Completing mission")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	log.Debug("Checking mission targets")

	// Get mission targets.
	targets, err := s.st.GetMissionTargets(ctx, mission)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			if errors.Is(err, context.DeadlineExceeded) {
				return models.Mission{}, models.ErrTimeoutExceeded
			}
			if errors.Is(err, pgx.ErrNoRows) {
				slog.Warn("Mission has no targets")
			}
		slog.Error("Failed to get mission targets", "err", err)
		return models.Mission{}, errors.New("failed to complete mission")
	}

	// Can't complete mission with pending targets.
	for _, t := range targets {
		if !t.Completed {
			log.Info("Mission has pending targets")
			return models.Mission{}, models.NewError(http.StatusUnprocessableEntity, "mission has pending targets")
		}
	}

	// Complete mission.
	completed, err := s.st.CompleteMission(ctx, mission)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Info("Mission not found")
			return models.Mission{}, models.ErrNotFound
		}
		slog.Error("Failed to complete mission", "err", err)
		return models.Mission{}, errors.New("failed to complete mission")
	}

	return sqlcMissionToModel(completed), nil
}

func (s Service) DeleteMission(ctx context.Context, id int32) error {
	log := slog.With(
		slog.String("op", "service.DeleteMission"),
		slog.Any("id", id),
	)

	log.Debug("Deleting mission")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get mission
	mission, err := s.st.GetMission(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("Mission not found")
			return models.ErrNotFound
		}
		slog.Error("Failed to get mission", "err", err)
		return errors.New("failed to delete mission")
	}

	// Can't delete an assigned mission
	if mission.Assignee.Valid {
		log.Info("Mission has assignee")
		return models.NewError(http.StatusUnprocessableEntity, "can't delete an assigned mission")
	}

	// Delete mission
	rows, err := s.st.DeleteMission(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.ErrTimeoutExceeded
		}
		slog.Error("Failed to delete mission", "err", err)
		return errors.New("failed to delete mission")
	}
	if rows == 0 {
		log.Debug("Mission not found")
		return models.ErrNotFound
	}

	log.Debug("Mission deleted")

	return nil
}

func sqlcTargetToModel(t postgres.Target) models.Target {
	return models.Target{
		ID:        t.ID,
		Name:      t.Name,
		Country:   t.Country,
		Notes:     t.Notes,
		Completed: t.Completed,
	}
}

func sqlcMissionToModel(t postgres.Mission) models.Mission {
	return models.Mission{
		ID:        t.ID,
		Assignee:  t.Assignee.Int32,
		Completed: t.Completed,
	}
}
