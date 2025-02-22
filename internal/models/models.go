package models

type Cat struct {
	Name              string `json:"name"`
	Breed             string `json:"breed"`
	YearsOfExperience int32  `json:"years_of_experience"`
	Salary            int32  `json:"salary"`
	ID                int32  `json:"id"`
}

type CreateCatRequest struct {
	Name              string `json:"name" validate:"required"`
	Breed             string `json:"breed" validate:"required"`
	YearsOfExperience int32  `json:"years_of_experience" validate:"required"`
	Salary            int32  `json:"salary" validate:"required"`
}

type UpdateCatSalaryRequest struct {
	Salary int32 `json:"salary" validate:"required"`
}

type Target struct {
	ID        int32  `json:"id"`
	Name      string `json:"name"`
	Country   string `json:"country"`
	Notes     string `json:"notes"`
	Completed bool   `json:"completed"`
}

type CreateTargetRequest struct {
	Name    string `json:"name" validate:"required"`
	Country string `json:"country" validate:"required"`
	Notes   string `json:"notes" validate:"required"`
}

type Mission struct {
	ID        int32    `json:"id"`
	Assignee  int32    `json:"assignee"`
	Targets   []Target `json:"targets"`
	Completed bool     `json:"completed"`
}

type CreateMissionRequest struct {
	Targets []CreateTargetRequest `json:"targets" validate:"required"`
}

type AssignCatRequest struct {
	Assignee int32 `json:"assignee" validate:"required"`
}

type UpdateTargetNotesRequest struct {
	Notes string `json:"notes" validate:"required"`
}
