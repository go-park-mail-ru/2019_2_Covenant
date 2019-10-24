package models

import "time"

type Session struct {
	ID      uint64
	UserID  uint64
	Expires time.Time
	Data    string
}
