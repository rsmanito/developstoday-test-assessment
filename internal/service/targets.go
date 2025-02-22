package service

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/rsmanito/developstoday-test-assessment/internal/models"
	"github.com/rsmanito/developstoday-test-assessment/internal/storage/postgres"
)

func (s Service) AddTarget(ctx context.Context, missionId int32, req models.CreateTargetRequest) (models.Mission, error) {
	log := slog.With(
		slog.String("op", "service.AddTarget"),
		slog.Any("missionId", missionId),
		slog.Any("target", req),
	)

	log.Debug("Adding target")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get mission.
	m, err := s.st.GetMission(ctx, missionId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("Mission not found")
			return models.Mission{}, models.ErrNotFound
		}
		slog.Error("Failed to get mission", "err", err)
		return models.Mission{}, errors.New("failed to add target")
	}

	// Get targets for mission
	targets, err := s.st.GetMissionTargets(ctx, missionId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			if errors.Is(err, context.DeadlineExceeded) {
				return models.Mission{}, models.ErrTimeoutExceeded
			}
			slog.Error("Failed to get mission", "err", err)
			return models.Mission{}, errors.New("failed to add target")
	}

	// Check if mission already has 3 targets.
	if len(targets) >= 3 {
		log.Info("Mission already has 3 targets")
		return models.Mission{}, models.NewError(http.StatusUnprocessableEntity, "mission has maximum targets (3)")
	}

	// Create new target.
	_, err = s.st.CreateTarget(ctx, postgres.CreateTargetParams{
		Mission: missionId,
		Name:    req.Name,
		Country: req.Country,
		Notes:   req.Notes,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Mission{}, models.ErrTimeoutExceeded
		}
		slog.Error("Failed to create target", "err", err)
		return models.Mission{}, errors.New("failed to add target")
	}

	// Get updated targets.
	targets, err = s.st.GetMissionTargets(ctx, missionId)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			if errors.Is(err, context.DeadlineExceeded) {
				return models.Mission{}, models.ErrTimeoutExceeded
			}
			slog.Error("Failed to get targets", "err", err)
			return models.Mission{}, errors.New("failed to add target")
	}

	// Populate mission with new targets.
	mission := sqlcMissionToModel(m)
	for _, t := range targets {
		mission.Targets = append(mission.Targets, sqlcTargetToModel(t))
	}

	log.Info("Added target")

	return mission, nil
}

func (s Service) DeleteTarget(ctx context.Context, targetId int32) error {
	log := slog.With(
		slog.String("op", "service.DeleteTarget"),
		slog.Any("targetId", targetId),
	)

	log.Debug("Deleting target")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := s.st.DeleteTarget(ctx, targetId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.ErrTimeoutExceeded
		}
		slog.Error("Failed to delete target", "err", err)
		return errors.New("failed to delete target")
	}
	if rows == 0 {
		log.Debug("Target not found")
		return models.ErrNotFound
	}

	log.Debug("Target deleted")

	return nil
}

func (s Service) UpdateTargetNotes(ctx context.Context, targetId int32, notes string) (models.Target, error) {
	log := slog.With(
		slog.String("op", "service.UpdateTargetNotes"),
		slog.Any("targetId", targetId),
		slog.Any("notes", notes),
	)

	log.Debug("Updating target notes")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Can't update notes to empty.
	if notes == "" {
		log.Debug("Notes are empty")
		return models.Target{}, models.NewError(http.StatusUnprocessableEntity, "notes can't be empty")
	}

	// Get target.
	target, err := s.st.GetTarget(ctx, targetId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Target{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("Target not found")
			return models.Target{}, models.ErrNotFound
		}
		slog.Error("Failed to delete target", "err", err)
		return models.Target{}, errors.New("failed to update target")
	}

	// Check if target is already completed.
	if target.Completed {
		log.Debug("Can't change notes of a completed target")
		return models.Target{}, models.NewError(http.StatusUnprocessableEntity, "Can't change notes of a completed target")
	}

	// Get mission.
	mission, err := s.st.GetMissionByTargetID(ctx, targetId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Target{}, models.ErrTimeoutExceeded
		}
		slog.Error("Failed to delete target", "err", err)
		return models.Target{}, errors.New("failed to update target")
	}

	// Check if mission is already completed.
	if mission.Completed {
		log.Debug("Can't change target notes of a completed mission")
		return models.Target{}, models.NewError(http.StatusUnprocessableEntity, "Can't change target notes of a completed mission")
	}
	res, err := s.st.UpdateTargetNotes(ctx, postgres.UpdateTargetNotesParams{
		ID:    targetId,
		Notes: notes,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Target{}, models.ErrTimeoutExceeded
		}
		slog.Error("Failed to delete target", "err", err)
		return models.Target{}, errors.New("failed to delete target")
	}

	return sqlcTargetToModel(res), nil
}

func (s Service) CompleteTarget(ctx context.Context, targetId int32) (models.Target, error) {
	log := slog.With(
		slog.String("op", "service.CompleteTarget"),
		slog.Any("targetId", targetId),
	)

	log.Debug("Completing target")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	// Get target.
	_, err := s.st.GetTarget(ctx, targetId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Target{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("Target not found")
			return models.Target{}, models.ErrNotFound
		}
		slog.Error("Failed to complete target", "err", err)
		return models.Target{}, errors.New("failed to complete target")
	}

	// Set target as completed.
	target, err := s.st.CompleteTarget(ctx, targetId)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Target{}, models.ErrTimeoutExceeded
		}
		slog.Error("Failed to complete target", "err", err)
		return models.Target{}, errors.New("failed to complete target")
	}

	return sqlcTargetToModel(target), nil
}
