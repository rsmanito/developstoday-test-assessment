-- name: GetAllCats :many
SELECT *
FROM CATS;

-- name: CreateCat :one
INSERT INTO cats (
  name, years_of_experience, breed, salary
) VALUES ( $1, $2, $3, $4)
RETURNING *;

-- name: GetCat :one
SELECT *
FROM cats 
WHERE id = $1
LIMIT 1;

-- name: UpdateCatSalary :one 
UPDATE cats
SET salary = $2
WHERE id = $1
RETURNING *;

-- name: DeleteCat :execrows
DELETE 
FROM cats
WHERE id = $1;
