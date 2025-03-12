-- name: CreateLog :exec
insert into logs (
    timestamp,
    caller_user_id,
    remote_addr,
    method,
    url,
    status,
    duration_ms
) values ($1, $2, $3, $4, $5, $6, $7);
