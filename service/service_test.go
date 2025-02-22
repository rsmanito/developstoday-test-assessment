package service

import (
	"context"
	"errors"
	"testing"

	"github.com/rsmanito/developstoday-test-assessment/models"
	"github.com/rsmanito/developstoday-test-assessment/storage/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockStorage struct {
	mock.Mock
}

func (m *MockStorage) GetAllCats(ctx context.Context) ([]postgres.Cat, error) {
	args := m.Called(ctx)
	return args.Get(0).([]postgres.Cat), args.Error(1)
}

func (m *MockStorage) CreateCat(ctx context.Context, params postgres.CreateCatParams) (postgres.Cat, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(postgres.Cat), args.Error(1)
}

func (m *MockStorage) GetCat(ctx context.Context, id int32) (postgres.Cat, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(postgres.Cat), args.Error(1)
}

func (m *MockStorage) UpdateCatSalary(ctx context.Context, params postgres.UpdateCatSalaryParams) (postgres.Cat, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(postgres.Cat), args.Error(1)
}

func (m *MockStorage) DeleteCat(ctx context.Context, id int32) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

/* ==============
   GET OPERATIONS
   ============== */

func TestGetAllCats(t *testing.T) {
	mockStorage := &MockStorage{}
	service := New(mockStorage)

	mockStorage.On("GetAllCats", mock.Anything).Return([]postgres.Cat{
		{ID: 1, Name: "Tom", Breed: "Siamese", YearsOfExperience: 5, Salary: 5000},
	}, nil)

	cats, err := service.GetAllCats(context.Background())
	assert.NoError(t, err)
	assert.Len(t, cats, 1)
	assert.Equal(t, "Tom", cats[0].Name)

	mockStorage.AssertExpectations(t)
}

func TestGetAllCats_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("GetAllCats", mock.Anything).Return([]postgres.Cat{}, errors.New("database error"))

	cats, err := service.GetAllCats(context.Background())
	assert.Error(t, err)
	assert.Equal(t, []models.Cat{}, cats)

	mockStorage.AssertExpectations(t)
}

func TestGetCat(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("GetCat", mock.Anything, int32(1)).Return(postgres.Cat{
		ID:                1,
		Name:              "Tom",
		Breed:             "Siamese",
		YearsOfExperience: 5,
		Salary:            5000,
	}, nil)

	cat, err := service.GetCat(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, "Tom", cat.Name)

	mockStorage.AssertExpectations(t)
}

func TestGetCat_NotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("GetCat", mock.Anything, int32(999)).Return(postgres.Cat{}, errors.New("not found"))

	cat, err := service.GetCat(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, models.Cat{}, cat)

	mockStorage.AssertExpectations(t)
}

/* =================
   CREATE OPERATIONS
   =================*/

func TestCreateCat(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("CreateCat", mock.Anything, postgres.CreateCatParams{
		Name:              "Tom",
		Breed:             "Siamese",
		YearsOfExperience: 5,
		Salary:            5000,
	}).Return(postgres.Cat{
		ID:                1,
		Name:              "Tom",
		Breed:             "Siamese",
		YearsOfExperience: 5,
		Salary:            5000,
	}, nil)

	cat, err := service.CreateCat(context.Background(), models.CreateCatRequest{
		Name:              "Tom",
		Breed:             "Siamese",
		YearsOfExperience: 5,
		Salary:            5000,
	})
	assert.NoError(t, err)
	assert.Equal(t, "Tom", cat.Name)

	mockStorage.AssertExpectations(t)
}

/* =================
   UPDATE OPERATIONS
   ================= */

func TestUpdateCatSalary(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("UpdateCatSalary", mock.Anything, postgres.UpdateCatSalaryParams{
		ID:     1,
		Salary: 6000,
	}).Return(postgres.Cat{
		ID:                1,
		Name:              "Tom",
		Breed:             "Siamese",
		YearsOfExperience: 5,
		Salary:            6000,
	}, nil)

	cat, err := service.UpdateCatSalary(context.Background(), models.UpdateCatSalaryRequest{
		Salary: 6000,
	}, 1)
	assert.NoError(t, err)
	assert.Equal(t, int32(6000), cat.Salary)

	mockStorage.AssertExpectations(t)
}

func TestUpdateCatSalary_InvalidID(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("UpdateCatSalary", mock.Anything, mock.Anything).Return(postgres.Cat{}, errors.New("not found"))

	_, err := service.UpdateCatSalary(context.Background(), models.UpdateCatSalaryRequest{Salary: 7000}, 999)
	assert.Error(t, err)
}

func TestUpdateCatSalary_NegativeSalary(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	_, err := service.UpdateCatSalary(context.Background(), models.UpdateCatSalaryRequest{Salary: -500}, 1)
	assert.Error(t, err)
}

/* =================
   DELETE OPERATIONS
   ================= */

func TestDeleteCat(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("DeleteCat", mock.Anything, int32(999)).Return(int64(1), nil)

	err := service.DeleteCat(context.Background(), 999)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}

func TestDeleteCat_NotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("DeleteCat", mock.Anything, int32(999)).Return(int64(0), nil)

	err := service.DeleteCat(context.Background(), 999)
	assert.Error(t, err)

	mockStorage.AssertExpectations(t)
}
