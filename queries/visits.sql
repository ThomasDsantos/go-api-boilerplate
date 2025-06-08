-- name: CountAllVisits :one
SELECT
  COUNT(*)
FROM
  visit;

-- name: CountVisitors :one
SELECT
  COUNT(DISTINCT(ip))
FROM
  visit
WHERE
  visited_at > $1;

-- name: InsertVisit :one
INSERT INTO visit (ip, name, visited_at)
VALUES (
  (sqlc.arg(ip)::varchar)::inet, sqlc.arg(name)::varchar, now()
) 
RETURNING *;
