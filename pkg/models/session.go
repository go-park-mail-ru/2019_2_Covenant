package models

import "time"

type Session struct {
	ID 		     uint
	UserID       uint
	Expiration   time.Time
	Data         string
}
