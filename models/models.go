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
	YearsOfExperience int32  `json:"years_of_experience" validate:"required,gte=0"`
	Salary            int32  `json:"salary" validate:"required,gte=0"`
}

type UpdateCatSalaryRequest struct {
	Salary int32 `json:"salary" validate:"required,gte=0"`
}
