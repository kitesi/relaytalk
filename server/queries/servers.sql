-- name: CreateServer :exec
insert into servers (owner_id, name, description) values ($1, $2, $3);
