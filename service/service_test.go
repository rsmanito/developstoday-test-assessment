package service

import (
	"context"
	"errors"
	"testing"
	"unicode/utf8"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

// Begin is a dummy implementation to satisfy the Storage interface.
func (m *MockStorage) Begin(ctx context.Context) (pgx.Tx, error) {
	return nil, nil
}

// CreateMission is a dummy implementation to satisfy the Storage interface.
func (m *MockStorage) CreateMission(ctx context.Context) (postgres.Mission, error) {
	args := m.Called(ctx)
	return args.Get(0).(postgres.Mission), args.Error(1)
}

// CreateTarget is a dummy implementation to satisfy the Storage interface.
func (m *MockStorage) CreateTarget(ctx context.Context, params postgres.CreateTargetParams) (postgres.Target, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(postgres.Target), args.Error(1)
}

// WithTx is a dummy implementation to satisfy the Storage interface.
func (m *MockStorage) WithTx(tx pgx.Tx) *postgres.Queries {
	return nil
}

func (m *MockStorage) GetAllMissions(ctx context.Context) ([]postgres.Mission, error) {
	args := m.Called(ctx)
	return args.Get(0).([]postgres.Mission), args.Error(1)
}

func (m *MockStorage) GetMission(ctx context.Context, id int32) (postgres.Mission, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(postgres.Mission), args.Error(1)
}

func (m *MockStorage) GetCatMission(ctx context.Context, id pgtype.Int4) (postgres.Mission, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(postgres.Mission), args.Error(1)
}

func (m *MockStorage) DeleteMission(ctx context.Context, id int32) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) AssignCat(ctx context.Context, params postgres.AssignCatParams) (postgres.Mission, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(postgres.Mission), args.Error(1)
}

func (m *MockStorage) CompleteMission(ctx context.Context, id int32) (postgres.Mission, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(postgres.Mission), args.Error(1)
}

func (m *MockStorage) GetMissionByTargetID(ctx context.Context, id int32) (postgres.GetMissionByTargetIDRow, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(postgres.GetMissionByTargetIDRow), args.Error(1)
}

func (m *MockStorage) GetTarget(ctx context.Context, id int32) (postgres.Target, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(postgres.Target), args.Error(1)
}

func (m *MockStorage) GetMissionTargets(ctx context.Context, id int32) ([]postgres.Target, error) {
	args := m.Called(ctx, id)
	return args.Get(0).([]postgres.Target), args.Error(1)
}

func (m *MockStorage) DeleteTarget(ctx context.Context, id int32) (int64, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(int64), args.Error(1)
}

func (m *MockStorage) UpdateTargetNotes(ctx context.Context, params postgres.UpdateTargetNotesParams) (postgres.Target, error) {
	args := m.Called(ctx, params)
	return args.Get(0).(postgres.Target), args.Error(1)
}

func (m *MockStorage) CompleteTarget(ctx context.Context, id int32) (postgres.Target, error) {
	args := m.Called(ctx, id)
	return args.Get(0).(postgres.Target), args.Error(1)
}


//-------------------------------------
// CATS TESTS
//-------------------------------------

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

//-------------------------------------
// MISSION TESTS
//-------------------------------------

func TestGetAllMissions_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockMissions := []postgres.Mission{
		{ID: 1, Assignee: pgtype.Int4{Int32: 2, Valid: true}, Completed: false},
		{ID: 2, Assignee: pgtype.Int4{Int32: 0, Valid: false}, Completed: true},
	}
	mockStorage.On("GetAllMissions", mock.Anything).Return(mockMissions, nil)

	missions, err := service.GetAllMissions(context.Background())
	assert.NoError(t, err)
	assert.Len(t, missions, 2)
	assert.Equal(t, int32(1), missions[0].ID)

	mockStorage.AssertExpectations(t)
}

func TestGetAllMissions_StorageError(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("GetAllMissions", mock.Anything).Return([]postgres.Mission{}, errors.New("db error"))

	missions, err := service.GetAllMissions(context.Background())
	assert.Error(t, err)
	assert.Len(t, missions, 0)

	mockStorage.AssertExpectations(t)
}

func TestGetMission_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	missionRecord := postgres.Mission{ID: 1, Assignee: pgtype.Int4{Int32: 2, Valid: true}, Completed: false}
	targetRecords := []postgres.Target{
		{ID: 10, Name: "Target1", Country: "CountryA", Notes: "Note", Completed: true},
		{ID: 11, Name: "Target2", Country: "CountryB", Notes: "Note", Completed: true},
	}

	mockStorage.On("GetMission", mock.Anything, int32(1)).Return(missionRecord, nil)
	mockStorage.On("GetMissionTargets", mock.Anything, int32(1)).Return(targetRecords, nil)

	mission, err := service.GetMission(context.Background(), 1)
	assert.NoError(t, err)
	assert.Equal(t, int32(1), mission.ID)
	assert.Len(t, mission.Targets, 2)

	mockStorage.AssertExpectations(t)
}

func TestGetMission_NotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	mockStorage.On("GetMission", mock.Anything, int32(999)).Return(postgres.Mission{}, errors.New("not found"))

	mission, err := service.GetMission(context.Background(), 999)
	assert.Error(t, err)
	assert.Equal(t, models.Mission{}, mission)

	mockStorage.AssertExpectations(t)
}

func TestCreateMission_InvalidTargetsCount(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	// Zero targets.
	reqZero := models.CreateMissionRequest{
		Targets: []models.CreateTargetRequest{},
	}
	_, err := service.CreateMission(context.Background(), reqZero)
	assert.Error(t, err)

	// Too many targets.
	reqTooMany := models.CreateMissionRequest{
		Targets: []models.CreateTargetRequest{
			{Name: "A", Country: "X", Notes: "N"},
			{Name: "B", Country: "Y", Notes: "N"},
			{Name: "C", Country: "Z", Notes: "N"},
			{Name: "D", Country: "W", Notes: "N"},
		},
	}
	_, err = service.CreateMission(context.Background(), reqTooMany)
	assert.Error(t, err)
}

func TestCreateMission_TargetNotesTooLong(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)

	// Target note longer than 256 characters.
	tooLong := ""
	for i := 0; i < 300; i++ {
		tooLong += "a"
	}

	assert.True(t, utf8.RuneCountInString(tooLong) > 256)

	req := models.CreateMissionRequest{
		Targets: []models.CreateTargetRequest{
			{Name: "Alpha", Country: "CountryX", Notes: tooLong},
		},
	}
	_, err := service.CreateMission(context.Background(), req)
	assert.Error(t, err)
}

func TestDeleteMission_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Get mission before deletion.
	missionRecord := postgres.Mission{ID: 1, Assignee: pgtype.Int4{Valid: false}}
	mockStorage.On("GetMission", mock.Anything, int32(1)).Return(missionRecord, nil)
	mockStorage.On("DeleteMission", mock.Anything, int32(1)).Return(int64(1), nil)

	err := service.DeleteMission(ctx, 1)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}

func TestDeleteMission_AssignedMission(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Cannot delete a mission that has an assignee.
	missionRecord := postgres.Mission{ID: 2, Assignee: pgtype.Int4{Int32: 10, Valid: true}}
	mockStorage.On("GetMission", mock.Anything, int32(2)).Return(missionRecord, nil)

	err := service.DeleteMission(ctx, 2)
	assert.Error(t, err)

	mockStorage.AssertExpectations(t)
}

//-------------------------------------
// TARGET TESTS
//-------------------------------------

func TestAddTarget_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Mission with no targets.
	missionRecord := postgres.Mission{ID: 1, Assignee: pgtype.Int4{Valid: false}, Completed: false}
	mockStorage.On("GetMission", mock.Anything, int32(1)).Return(missionRecord, nil)
	mockStorage.On("GetMissionTargets", mock.Anything, int32(1)).Return([]postgres.Target{}, nil).Once()

	newTargetParams := postgres.CreateTargetParams{
		Mission: 1,
		Name:    "TargetX",
		Country: "CountryX",
		Notes:   "Note",
	}
	newTarget := postgres.Target{ID: 100, Name: "TargetX", Country: "CountryX", Notes: "Note", Completed: false}
	mockStorage.On("CreateTarget", mock.Anything, newTargetParams).Return(newTarget, nil)

	// Get mission after adding target.
	mockStorage.On("GetMissionTargets", mock.Anything, int32(1)).Return([]postgres.Target{newTarget}, nil).Once()

	req := models.CreateTargetRequest{
		Name:    "TargetX",
		Country: "CountryX",
		Notes:   "Note",
	}
	mission, err := service.AddTarget(ctx, 1, req)
	assert.NoError(t, err)
	assert.Len(t, mission.Targets, 1)
	assert.Equal(t, "TargetX", mission.Targets[0].Name)

	mockStorage.AssertExpectations(t)
}

func TestAddTarget_TooManyTargets(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Mission with three existing targets.
	missionRecord := postgres.Mission{ID: 1, Assignee: pgtype.Int4{Valid: false}, Completed: false}
	targets := []postgres.Target{
		{ID: 1}, {ID: 2}, {ID: 3},
	}
	mockStorage.On("GetMission", mock.Anything, int32(1)).Return(missionRecord, nil)
	mockStorage.On("GetMissionTargets", mock.Anything, int32(1)).Return(targets, nil)

	req := models.CreateTargetRequest{
		Name:    "ExtraTarget",
		Country: "CountryY",
		Notes:   "Note",
	}
	_, err := service.AddTarget(ctx, 1, req)
	assert.Error(t, err)

	mockStorage.AssertExpectations(t)
}

func TestDeleteTarget_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	mockStorage.On("DeleteTarget", mock.Anything, int32(100)).Return(int64(1), nil)

	err := service.DeleteTarget(ctx, 100)
	assert.NoError(t, err)

	mockStorage.AssertExpectations(t)
}

func TestDeleteTarget_NotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	mockStorage.On("DeleteTarget", mock.Anything, int32(200)).Return(int64(0), nil)

	err := service.DeleteTarget(ctx, 200)
	assert.Error(t, err)

	mockStorage.AssertExpectations(t)
}

func TestUpdateTargetNotes_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Pending target.
	targetRecord := postgres.Target{ID: 150, Name: "Target150", Country: "Country150", Notes: "Old", Completed: false}
	mockStorage.On("GetTarget", mock.Anything, int32(150)).Return(targetRecord, nil)

	// Pending mission.
	missionRow := postgres.GetMissionByTargetIDRow{MissionID: 10, Completed: false}
	mockStorage.On("GetMissionByTargetID", mock.Anything, int32(150)).Return(missionRow, nil)

	updateParams := postgres.UpdateTargetNotesParams{
		ID:    150,
		Notes: "New Notes",
	}
	updatedTarget := postgres.Target{ID: 150, Name: "Target150", Country: "Country150", Notes: "New Notes", Completed: false}
	mockStorage.On("UpdateTargetNotes", mock.Anything, updateParams).Return(updatedTarget, nil)

	target, err := service.UpdateTargetNotes(ctx, 150, "New Notes")
	assert.NoError(t, err)
	assert.Equal(t, "New Notes", target.Notes)

	mockStorage.AssertExpectations(t)
}

func TestUpdateTargetNotes_Empty(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Can't update target notes to empty.
	target, err := service.UpdateTargetNotes(ctx, 150, "")
	assert.Error(t, err)
	assert.Equal(t, models.Target{}, target)
}

func TestCompleteTarget_Success(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Pending target.
	targetRecord := postgres.Target{ID: 170, Completed: false}
	mockStorage.On("GetTarget", mock.Anything, int32(170)).Return(targetRecord, nil)

	// Pending mission.
	completedTarget := postgres.Target{ID: 170, Completed: true}
	mockStorage.On("CompleteTarget", mock.Anything, int32(170)).Return(completedTarget, nil)

	target, err := service.CompleteTarget(ctx, 170)
	assert.NoError(t, err)
	assert.True(t, target.Completed)

	mockStorage.AssertExpectations(t)
}

func TestCompleteTarget_NotFound(t *testing.T) {
	mockStorage := new(MockStorage)
	service := New(mockStorage)
	ctx := context.Background()

	// Target not found.
	mockStorage.On("GetTarget", mock.Anything, int32(180)).Return(postgres.Target{}, errors.New("not found"))

	target, err := service.CompleteTarget(ctx, 180)
	assert.Error(t, err)
	assert.Equal(t, models.Target{}, target)

	mockStorage.AssertExpectations(t)
}