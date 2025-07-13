-- name: CreateChannel :exec
insert into channels (server_id, owner_id, name, description) values ($1, $2, $3, $4);
