-- name: QueryGetAlluser :many
SELECT * FROM user_details;

-- name: QueryAddUser :exec
INSERT INTO  user_details (email,user_password,created,updated) values ($1,$2,$3,$4);

-- name: QueryCheckUserByEmail :one
SELECT * from user_details WHERE email=$1;

-- name: QueryGetUserById :one
SELECT * from user_details WHERE id=$1;




