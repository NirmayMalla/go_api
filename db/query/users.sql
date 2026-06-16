-- name: CreateUser :execresult
INSERT INTO users(name, dob)
VALUES(?, ?);

-- name: GetUser :one
SELECT * 
FROM users
WHERE id=?;

-- name: ListUsers :many
SELECT *
FROM users;

-- name: UpdateUser :execresult
UPDATE users
SET name=?, dob=?
WHERE id=?;

-- name: DeleteUser :execresult
DELETE FROM users
WHERE id=?
