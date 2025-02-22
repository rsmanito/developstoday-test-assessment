package service

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v3/client"
	"github.com/jackc/pgx/v5"
	"github.com/rsmanito/developstoday-test-assessment/models"
	"github.com/rsmanito/developstoday-test-assessment/storage/postgres"
)

func (s Service) GetAllCats(ctx context.Context) ([]models.Cat, error) {
	log := slog.With(
		slog.String("op", "service.GetCats"),
	)
	log.Debug("Fetching all cats")

	ctx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	cats, err := s.st.GetAllCats(ctx)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return make([]models.Cat, 0), models.ErrTimeoutExceeded
		} else {
			return make([]models.Cat, 0), errors.New("failed to fetch cats")
		}
	}

	mappedCats := make([]models.Cat, len(cats))
	for i, cat := range cats {
		mappedCats[i] = sqlcCatToModel(cat)
	}

	log.Debug("Got cats", "cats", mappedCats)

	return mappedCats, nil
}

func (s Service) CreateCat(ctx context.Context, req models.CreateCatRequest) (models.Cat, error) {
	log := slog.With(
		slog.String("op", "service.CreateCat"),
		slog.Any("req", req),
	)
	log.Debug("Creating a cat")

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Validate breed
	breedReqCtx, breedCancel := context.WithTimeout(ctx, 2*time.Second)
	defer breedCancel()

	resChan := make(chan bool, 1)
	errChan := make(chan error, 1)

	go breedExists(req.Breed, resChan, errChan)

	select {
	case res := <-resChan:
		if !res {
			log.Warn("Unknown breed")
			return models.Cat{}, models.NewError(http.StatusUnprocessableEntity, "Unknown breed")
		}
	case err := <-errChan:
		log.Error("Failed to validate cat breed", "err", err)
		return models.Cat{}, errors.New("failed to create a cat")
	case <-breedReqCtx.Done():
		log.Debug("Validate cat breed timeout")
		return models.Cat{}, errors.New("failed to create a cat")
	}

	// Create the cat in the database
	dbCtx, cancel := context.WithTimeout(ctx, 500*time.Millisecond)
	defer cancel()

	res, err := s.st.CreateCat(dbCtx, postgres.CreateCatParams{
		Name:              req.Name,
		Breed:             req.Breed,
		YearsOfExperience: req.YearsOfExperience,
		Salary:            req.Salary,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Cat{}, models.ErrTimeoutExceeded
		}
		log.Error("Failed to save cat", "err", err)
		return models.Cat{}, errors.New("failed to create a cat")
	}

	return sqlcCatToModel(res), nil
}

func (s Service) GetCat(ctx context.Context, id int32) (models.Cat, error) {
	log := slog.With(
		slog.String("op", "service.GetCat"),
		slog.Any("id", id),
	)

	log.Debug("Fetching cat")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := s.st.GetCat(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Cat{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			log.Debug("Cat not found")
			return models.Cat{}, models.ErrNotFound
		}
		log.Error("Failed to get cat", "err", err)
		return models.Cat{}, errors.New("Failed to get cat")
	}

	return sqlcCatToModel(res), nil
}

func (s Service) UpdateCatSalary(ctx context.Context, req models.UpdateCatSalaryRequest, id int32) (models.Cat, error) {
	if req.Salary < 0 {
		return models.Cat{}, models.NewError(http.StatusUnprocessableEntity, "salary must be greater than or equal to 0")
	}
	log := slog.With(
		slog.String("op", "service.UpdateCatSalary"),
		slog.Any("id", id),
	)

	log.Debug("Updating cat salary")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	res, err := s.st.UpdateCatSalary(ctx, postgres.UpdateCatSalaryParams{
		ID:     id,
		Salary: req.Salary,
	})
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.Cat{}, models.ErrTimeoutExceeded
		}
		if errors.Is(err, pgx.ErrNoRows) {
			return models.Cat{}, models.ErrNotFound
		}
		log.Error("Failed to update salary", "err", err)
		return models.Cat{}, errors.New("failed to update salary")
	}

	log.Debug("Updated cat salary")

	return sqlcCatToModel(res), nil
}

func (s Service) DeleteCat(ctx context.Context, id int32) error {
	log := slog.With(
		slog.String("op", "service.UpdateCatSalary"),
		slog.Any("id", id),
	)

	log.Debug("Deleting cat")

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := s.st.DeleteCat(ctx, id)
	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return models.ErrTimeoutExceeded
		}
		slog.Error("Failed to delete cat", "err", err)
		return errors.New("failed to delete cat")
	}
	if rows == 0 {
		return models.ErrNotFound
	}

	log.Debug("Deleted cat")

	return nil
}

func sqlcCatToModel(c postgres.Cat) models.Cat {
	return models.Cat{
		ID:                c.ID,
		Name:              c.Name,
		Breed:             c.Breed,
		YearsOfExperience: c.YearsOfExperience,
		Salary:            c.Salary,
	}
}

func breedExists(breedName string, resChan chan<- bool, errChan chan<- error) {
	cc := client.New()

	// Get cat breeds.
	resp, err := cc.Get("https://api.thecatapi.com/v1/breeds")
	if err != nil {
		errChan <- err
		return
	}

	type Breed struct {
		Name string `json:"name"`
	}

	var breeds []Breed
	if err := json.Unmarshal(resp.Body(), &breeds); err != nil {
		errChan <- err
		return
	}

	// Check if the breed exists in the list.
	for _, breed := range breeds {
		if strings.EqualFold(breed.Name, breedName) {
			resChan <- true
			return
		}
	}
	resChan <- false
}
