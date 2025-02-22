-- name: GetAllCats :many
SELECT *
FROM cats;

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

-- name: GetAllMissions :many
SELECT *
FROM missions;

-- name: CreateMission :one
INSERT INTO missions
DEFAULT VALUES
RETURNING *;

-- name: GetMission :one
SELECT *
FROM missions
WHERE missions.id = $1;

-- name: DeleteMission :execrows
DELETE
FROM missions
WHERE id = $1;

-- name: AssignCat :one
UPDATE missions
SET assignee = $2
WHERE id = $1
RETURNING *;

-- name: GetCatMission :one
SELECT *
FROM missions
WHERE assignee = $1
LIMIT 1;

-- name: CompleteMission :one
UPDATE missions
SET completed = true
WHERE id = $1
RETURNING *;

-- name: CreateTarget :one
INSERT INTO targets (
  mission, name, country, notes
) VALUES ( $1, $2, $3, $4 )
RETURNING *;

-- name: GetMissionTargets :many
SELECT *
FROM targets
WHERE mission = $1;

-- name: DeleteTarget :execrows
DELETE 
FROM targets
WHERE id = $1;

-- name: UpdateTargetNotes :one
UPDATE targets
SET notes = $2
WHERE id = $1
RETURNING *;

-- name: GetMissionByTargetID :one
SELECT 
    m.id AS mission_id,
    m.assignee,
    m.completed
FROM missions m
JOIN targets t ON m.id = t.mission
WHERE t.id = $1;

-- name: GetTarget :one
SELECT *
FROM targets
WHERE id = $1
LIMIT 1;

-- name: CompleteTarget :one
UPDATE targets
SET completed = true
WHERE id = $1
RETURNING *;
