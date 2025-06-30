-- name: CreateUser :exec
insert into users (username, password_hash, email) values ($1, $2, $3);


-- name: GetUserByUsername :one
select id, username, password_hash from users where username = $1;

-- name: GetUserByEmail :one
select id, username, password_hash from users where email = $1;
