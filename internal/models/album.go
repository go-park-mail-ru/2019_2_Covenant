package models

import "time"

type Album struct {
	ID       uint64    `json:"-"`
	ArtistID string    `json:"artist_id"`
	Name     string    `json:"name"`
	Photo    string    `json:"photo"`
	Year     time.Time `json:"year"`
}
