// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0
// source: visits.sql

package store

import (
	"context"
	"time"
)

const countAllVisits = `-- name: CountAllVisits :one
SELECT
  COUNT(*)
FROM
  visit
`

func (q *Queries) CountAllVisits(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, countAllVisits)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const countVisitors = `-- name: CountVisitors :one
SELECT
  COUNT(DISTINCT(ip))
FROM
  visit
WHERE
  visited_at > $1
`

func (q *Queries) CountVisitors(ctx context.Context, visitedAt time.Time) (int64, error) {
	row := q.db.QueryRow(ctx, countVisitors, visitedAt)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const insertVisit = `-- name: InsertVisit :one
INSERT INTO visit (ip, name, visited_at)
VALUES (
  ($1::varchar)::inet, $2::varchar, now()
) 
RETURNING ip, name, visited_at
`

type InsertVisitParams struct {
	Ip   string
	Name string
}

func (q *Queries) InsertVisit(ctx context.Context, arg InsertVisitParams) (Visit, error) {
	row := q.db.QueryRow(ctx, insertVisit, arg.Ip, arg.Name)
	var i Visit
	err := row.Scan(&i.Ip, &i.Name, &i.VisitedAt)
	return i, err
}
