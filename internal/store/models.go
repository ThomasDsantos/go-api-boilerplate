// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.29.0

package store

import (
	"net"
	"time"
)

type Visit struct {
	Ip        net.IP
	Name      string
	VisitedAt time.Time
}
