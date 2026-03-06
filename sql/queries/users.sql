-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES(
    $1,
    $2,
    $3,
    $4,
    $5
)
RETURNING *;




-- name: GetUser :one
SELECT * FROM users
WHERE id = $1;

-- name: ResetUsers :exec
DELETE FROM users;

-- name: GetUsers :many
SELECT id FROM users;

-- name: GetUserFromEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: ChangeUsersEmailAndPassword :exec
UPDATE users
SET email = $1,
hashed_password =  $2
WHERE id = $3;


-- name: SetUsersChirpyRed :exec
UPDATE users
SET is_chirpy_red = $1
WHERE id = $2;
