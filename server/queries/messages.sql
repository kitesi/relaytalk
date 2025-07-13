-- name: CreateMessage :exec
insert into messages (user_id, content, channel_id) values ($1, $2, $3);
